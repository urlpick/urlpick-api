package url

import "time"

type URL struct {
	ID          uint      `json:"id"`
	Hash        string    `json:"hash"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}
