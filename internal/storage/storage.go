package storage

import "github.com/yourusername/url-shortener/internal/models"

type Storage interface {
	Save(url models.URL) error
	Get(shortCode string) (models.URL, bool)
	IncrementClicks(shortCode string) error
	Exists(shortCode string) bool
	SaveToFile() error
	GetAll() map[string]models.URL
}
