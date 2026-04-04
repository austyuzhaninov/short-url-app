package models

import "time"

type URL struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	UserID      string    `json:"user_id"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShortenRequest struct {
	URL    string `json:"url"`
	UserID string `json:"user_id"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

type StatsResponse struct {
	OriginalURL string    `json:"original_url"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}
