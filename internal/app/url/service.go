package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/urlpick/urlpick-api/internal/pkg/config"
	"github.com/urlpick/urlpick-api/internal/pkg/utils/errors"
	"github.com/urlpick/urlpick-api/internal/pkg/utils/turnstile"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateShortURL(ctx context.Context, req CreateURLRequest) (*CreateURLResponse, error) {
	if !turnstile.VerifyTurnstileToken(req.TurnstileToken) {
		return nil, errors.Unauthorized("invalid turnstile token")
	}

	hash, err := generateHash()
	if err != nil {
		return nil, err
	}

	url := &URL{
		Hash:        hash,
		OriginalURL: req.URL,
	}

	if err := s.repo.Save(ctx, url); err != nil {
		return nil, err
	}

	return &CreateURLResponse{
		ShortURL:    config.AppConfig.BaseURL + "/" + url.Hash,
		OriginalURL: url.OriginalURL,
		CreatedAt:   url.CreatedAt,
	}, nil
}

func (s *Service) GetURL(ctx context.Context, hash string) (*URLResponse, error) {
	url, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	return &URLResponse{
		OriginalURL: url.OriginalURL,
	}, nil
}

func generateHash() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:8], nil
}
