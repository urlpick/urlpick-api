package mysql

import (
	"context"
	"database/sql"

	"github.com/urlpick/urlpick-api/internal/app/url"
	"github.com/urlpick/urlpick-api/internal/pkg/utils/errors"
)

type urlRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) url.Repository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Save(ctx context.Context, url *url.URL) error {
	query := `
        INSERT INTO urls (hash, original_url)
        VALUES (?, ?)
    `
	result, err := r.db.ExecContext(ctx, query, url.Hash, url.OriginalURL)
	if err != nil {
		return errors.Internal("Failed to save URL")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Internal("Failed to get URL ID")
	}

	url.ID = uint(id)
	return nil
}

func (r *urlRepository) FindByHash(ctx context.Context, hash string) (*url.URL, error) {
	query := `
        SELECT id, hash, original_url, created_at, last_visit
        FROM urls WHERE hash = ?
    `
	var u url.URL
	err := r.db.QueryRowContext(ctx, query, hash).Scan(
		&u.ID,
		&u.Hash,
		&u.OriginalURL,
		&u.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NotFound("URL not found")
	}
	if err != nil {
		return nil, errors.Internal("Failed to get URL")
	}

	return &u, nil
}
