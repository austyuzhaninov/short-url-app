package entity

import "time"

// URL - сущность для хранения в файле/БД
type URL struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	UserID      string    `json:"user_id"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}
