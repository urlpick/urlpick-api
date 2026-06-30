package url

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/urlpick/urlpick-api/internal/pkg/utils/errors"
)

type fakeRepo struct {
	saveErrs  []error // returned in order, then nil
	saveCalls int
}

func (f *fakeRepo) Save(ctx context.Context, u *URL) error {
	var err error
	if f.saveCalls < len(f.saveErrs) {
		err = f.saveErrs[f.saveCalls]
	}
	f.saveCalls++
	return err
}

func (f *fakeRepo) FindByHash(ctx context.Context, hash string) (*URL, error) {
	return nil, nil
}

func passVerify(ctx context.Context, token string) bool { return true }

func newTestService(repo Repository, verify func(context.Context, string) bool) *Service {
	return &Service{repo: repo, verify: verify}
}

func TestCreateShortURL_DuplicateThenSuccess(t *testing.T) {
	tests := []struct {
		name      string
		dupCount  int
		wantCalls int
	}{
		{"no collision", 0, 1},
		{"one collision", 1, 2},
		{"four collisions", 4, 5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := make([]error, tc.dupCount)
			for i := range errs {
				errs[i] = ErrDuplicateHash
			}
			repo := &fakeRepo{saveErrs: errs}
			s := newTestService(repo, passVerify)

			resp, err := s.CreateShortURL(context.Background(), CreateURLRequest{
				URL:            "https://example.com",
				TurnstileToken: "tok",
			})
			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected response, got nil")
			}
			if repo.saveCalls != tc.wantCalls {
				t.Fatalf("expected %d save calls, got %d", tc.wantCalls, repo.saveCalls)
			}
		})
	}
}

func TestCreateShortURL_AlwaysDuplicate(t *testing.T) {
	errs := make([]error, maxSaveAttempts)
	for i := range errs {
		errs[i] = ErrDuplicateHash
	}
	repo := &fakeRepo{saveErrs: errs}
	s := newTestService(repo, passVerify)

	_, err := s.CreateShortURL(context.Background(), CreateURLRequest{
		URL:            "https://example.com",
		TurnstileToken: "tok",
	})
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if repo.saveCalls != maxSaveAttempts {
		t.Fatalf("expected %d save attempts, got %d", maxSaveAttempts, repo.saveCalls)
	}
}

func TestCreateShortURL_NonDuplicateErrorImmediate(t *testing.T) {
	sentinel := stderrors.New("db down")
	repo := &fakeRepo{saveErrs: []error{sentinel}}
	s := newTestService(repo, passVerify)

	_, err := s.CreateShortURL(context.Background(), CreateURLRequest{
		URL:            "https://example.com",
		TurnstileToken: "tok",
	})
	if !stderrors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if repo.saveCalls != 1 {
		t.Fatalf("expected 1 save call, got %d", repo.saveCalls)
	}
}

func TestCreateShortURL_InvalidTurnstile(t *testing.T) {
	repo := &fakeRepo{}
	s := newTestService(repo, func(context.Context, string) bool { return false })

	_, err := s.CreateShortURL(context.Background(), CreateURLRequest{
		URL:            "https://example.com",
		TurnstileToken: "tok",
	})
	var ce errors.CustomError
	if !stderrors.As(err, &ce) || ce.Status != 401 {
		t.Fatalf("expected 401 unauthorized, got %v", err)
	}
	if repo.saveCalls != 0 {
		t.Fatalf("expected no save calls, got %d", repo.saveCalls)
	}
}

func TestCreateShortURL_SchemeValidation(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"ftp", "ftp://example.com", true},
		{"javascript", "javascript:alert(1)", true},
		{"empty", "", true},
		{"no host", "http://", true},
		{"http", "http://example.com", false},
		{"https", "https://example.com", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{}
			s := newTestService(repo, passVerify)

			_, err := s.CreateShortURL(context.Background(), CreateURLRequest{
				URL:            tc.url,
				TurnstileToken: "tok",
			})
			if tc.wantErr {
				var ce errors.CustomError
				if !stderrors.As(err, &ce) || ce.Status != 400 {
					t.Fatalf("expected 400 bad request, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected success, got %v", err)
			}
		})
	}
}

func TestGenerateHash(t *testing.T) {
	const iterations = 10000
	seen := make(map[string]struct{}, iterations)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	for i := 0; i < iterations; i++ {
		h, err := generateHash()
		if err != nil {
			t.Fatalf("generateHash error: %v", err)
		}
		if len(h) != hashLength {
			t.Fatalf("expected length %d, got %d (%q)", hashLength, len(h), h)
		}
		for _, c := range h {
			if !containsRune(charset, c) {
				t.Fatalf("hash %q contains non-base64url char %q", h, c)
			}
		}
		if _, dup := seen[h]; dup {
			t.Fatalf("duplicate hash generated: %q", h)
		}
		seen[h] = struct{}{}
	}
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
