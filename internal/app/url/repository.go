package url

import "context"

type Repository interface {
	Save(ctx context.Context, url *URL) error
	FindByHash(ctx context.Context, hash string) (*URL, error)
}
