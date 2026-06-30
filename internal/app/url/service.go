package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	stderrors "errors"
	neturl "net/url"

	"github.com/urlpick/urlpick-api/internal/pkg/config"
	"github.com/urlpick/urlpick-api/internal/pkg/utils/errors"
	"github.com/urlpick/urlpick-api/internal/pkg/utils/turnstile"
)

// hashLength is the number of url-safe characters in a short hash (matches VARCHAR(8)).
const hashLength = 8

// maxSaveAttempts bounds hash regeneration on collision.
const maxSaveAttempts = 5

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateShortURL(ctx context.Context, req CreateURLRequest) (*CreateURLResponse, error) {
	if !turnstile.VerifyTurnstileToken(ctx, req.TurnstileToken) {
		return nil, errors.Unauthorized("invalid turnstile token")
	}

	parsed, err := neturl.Parse(req.URL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, errors.BadRequest("only http/https URLs are allowed")
	}

	var u *URL
	for attempt := 0; attempt < maxSaveAttempts; attempt++ {
		hash, err := generateHash()
		if err != nil {
			return nil, err
		}

		u = &URL{
			Hash:        hash,
			OriginalURL: req.URL,
		}

		err = s.repo.Save(ctx, u)
		if err == nil {
			return &CreateURLResponse{
				ShortURL:    config.AppConfig.BaseURL + "/" + u.Hash,
				OriginalURL: u.OriginalURL,
				CreatedAt:   u.CreatedAt,
			}, nil
		}
		if !stderrors.Is(err, ErrDuplicateHash) {
			return nil, err
		}
	}

	return nil, errors.Internal("failed to generate a unique short URL")
}

func (s *Service) GetURL(ctx context.Context, hash string) (*URLResponse, error) {
	u, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	return &URLResponse{
		OriginalURL: u.OriginalURL,
	}, nil
}

func generateHash() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:hashLength], nil
}
