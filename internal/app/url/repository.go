package url

import (
	"context"
	"errors"
)

// ErrDuplicateHash is returned by Repository.Save when a hash collides with an existing row.
var ErrDuplicateHash = errors.New("duplicate hash")

type Repository interface {
	Save(ctx context.Context, url *URL) error
	FindByHash(ctx context.Context, hash string) (*URL, error)
}
