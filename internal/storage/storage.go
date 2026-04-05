package storage

import (
	"short-url-app/internal/models/entity"
)

type Storage interface {
	Save(url entity.URL) error
	Get(shortCode string) (entity.URL, bool)
	IncrementClicks(shortCode string) error
	Exists(shortCode string) bool
	SaveToFile() error
	GetAll() map[string]entity.URL
}
