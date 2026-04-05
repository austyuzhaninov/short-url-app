package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"short-url-app/internal/models/entity"
	"sync"
)

type MemoryStorage struct {
	mu       sync.RWMutex
	urls     map[string]entity.URL
	filePath string
}

func NewMemoryStorage(filePath string) (*MemoryStorage, error) {
	storage := &MemoryStorage{
		urls:     make(map[string]entity.URL),
		filePath: filePath,
	}

	// Загрузка из файла, если существует
	if err := storage.loadFromFile(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load storage: %w", err)
		}
		// Файл не существует — создаём пустой
		if err := storage.SaveToFile(); err != nil {
			return nil, fmt.Errorf("failed to create storage file: %w", err)
		}
	}

	return storage, nil
}

func (s *MemoryStorage) Save(url entity.URL) error {
	s.mu.Lock()
	s.urls[url.ShortCode] = url
	s.mu.Unlock() // ← освобождаем ДО записи в файл

	return s.SaveToFile() // ← RLock() сможет захватиться
}

func (s *MemoryStorage) Get(shortCode string) (entity.URL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urls[shortCode]
	return url, exists
}

func (s *MemoryStorage) IncrementClicks(shortCode string) error {
	// 1. Захватываем Lock, обновляем данные
	s.mu.Lock()
	url, exists := s.urls[shortCode]
	if !exists {
		s.mu.Unlock()
		return nil
	}
	url.Clicks++
	s.urls[shortCode] = url
	s.mu.Unlock() // ← освобождаем ДО вызова SaveToFile

	// 2. Безопасно сохраняем в файл
	return s.SaveToFile()
}

func (s *MemoryStorage) Exists(shortCode string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.urls[shortCode]
	return exists
}

func (s *MemoryStorage) GetAll() map[string]entity.URL {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]entity.URL, len(s.urls))
	for k, v := range s.urls {
		result[k] = v
	}
	return result
}

func (s *MemoryStorage) SaveToFile() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.urls, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal storage: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write storage file: %w", err)
	}

	return nil
}

func (s *MemoryStorage) loadFromFile() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var urls map[string]entity.URL
	if err := json.Unmarshal(data, &urls); err != nil {
		return fmt.Errorf("failed to unmarshal storage: %w", err)
	}

	s.urls = urls
	return nil
}
