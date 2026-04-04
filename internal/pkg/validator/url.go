package validator

import (
	"net/url"
	"strings"
)

func IsValidURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	// Проверка на http:// или https://
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return false
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return parsed.Host != ""
}
