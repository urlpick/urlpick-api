package url

import "time"

type URL struct {
	ID          uint
	Hash        string
	OriginalURL string
	CreatedAt   time.Time
}
