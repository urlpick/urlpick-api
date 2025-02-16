package url

import "time"

type CreateURLRequest struct {
	URL            string `json:"url" binding:"required,url"`
	TurnstileToken string `json:"turnstile_token" binding:"required"`
}

type CreateURLResponse struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type URLResponse struct {
	OriginalURL string `json:"original_url"`
}
