package dto

import "time"

// ShortenRequest - запрос на создание короткой ссылки
type ShortenRequest struct {
	URL    string `json:"url" validate:"required,url"`
	UserID string `json:"user_id" validate:"required"`
}

// ShortenResponse - ответ с короткой ссылкой
type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

// StatsResponse - ответ со статистикой
type StatsResponse struct {
	OriginalURL string    `json:"original_url"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}
