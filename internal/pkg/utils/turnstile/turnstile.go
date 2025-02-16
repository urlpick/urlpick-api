package turnstile

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/urlpick/urlpick-api/internal/pkg/config"
)

const TURNSTILE_VERIFY_URL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

func VerifyTurnstileToken(token string) bool {
	if token == "" {
		return false
	}

	data := url.Values{}
	data.Set("secret", config.AppConfig.TurnstileSecretKey)
	data.Set("response", token)

	resp, err := http.Post(
		TURNSTILE_VERIFY_URL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
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
