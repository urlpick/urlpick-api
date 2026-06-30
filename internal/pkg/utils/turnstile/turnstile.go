package turnstile

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/urlpick/urlpick-api/internal/pkg/config"
)

const TURNSTILE_VERIFY_URL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

var httpClient = &http.Client{Timeout: 5 * time.Second}

func VerifyTurnstileToken(ctx context.Context, token string) bool {
	if token == "" {
		return false
	}

	data := url.Values{}
	data.Set("secret", config.AppConfig.TurnstileSecretKey)
	data.Set("response", token)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		TURNSTILE_VERIFY_URL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	type VerifyResponse struct {
		Success bool `json:"success"`
	}

	var result VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	return result.Success
}
