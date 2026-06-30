package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	driver "github.com/go-sql-driver/mysql"
	"github.com/urlpick/urlpick-api/internal/app/url"
	apperrors "github.com/urlpick/urlpick-api/internal/pkg/utils/errors"
)

const saveTimeout = 3 * time.Second

type urlRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) url.Repository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Save(ctx context.Context, u *url.URL) error {
	ctx, cancel := context.WithTimeout(ctx, saveTimeout)
	defer cancel()

	query := `
        INSERT INTO urls (hash, original_url, created_at)
        VALUES (?, ?, ?)
    `
	result, err := r.db.ExecContext(ctx, query, u.Hash, u.OriginalURL, u.CreatedAt)
	if err != nil {
		var mysqlErr *driver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return url.ErrDuplicateHash
		}
		return apperrors.Internal("Failed to save URL")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return apperrors.Internal("Failed to get URL ID")
	}

	u.ID = uint(id)
	return nil
}

func (r *urlRepository) FindByHash(ctx context.Context, hash string) (*url.URL, error) {
	ctx, cancel := context.WithTimeout(ctx, saveTimeout)
	defer cancel()

	query := `
        SELECT id, hash, original_url, created_at
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
		return nil, apperrors.NotFound("URL not found")
	}
	if err != nil {
		return nil, apperrors.Internal("Failed to get URL")
	}

	return &u, nil
}
