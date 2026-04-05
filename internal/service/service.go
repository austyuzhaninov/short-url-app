package service

import (
	"errors"
	"fmt"
	"short-url-app/internal/models/entity"
	"short-url-app/internal/pkg/generator"
	"short-url-app/internal/storage"
	"time"
)

type URLServiceInterface interface {
	ShortenURL(originalURL, userID string) (string, string, error)
	GetOriginalURL(shortCode string) (string, error)
	GetStats(shortCode string) (originalURL string, clicks int, createdAt time.Time, err error)
}

type URLService struct {
	storage storage.Storage
	baseURL string
}

func New(storage storage.Storage, baseURL string) *URLService {
	return &URLService{
		storage: storage,
		baseURL: baseURL,
	}
}

func (s *URLService) ShortenURL(originalURL, userID string) (shortCode string, shortURL string, err error) {
	// Генерация уникального кода
	for {
		shortCode = generator.GenerateShortCode()
		if !s.storage.Exists(shortCode) {
			break
		}
	}

	// Создаём entity для хранения
	url := entity.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		UserID:      userID,
		Clicks:      0,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.storage.Save(url); err != nil {
		return "", "", err
	}

	return shortCode, fmt.Sprintf("%s/%s", s.baseURL, shortCode), nil
}

func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	url, exists := s.storage.Get(shortCode)
	if !exists {
		return "", errors.New("short code not found")
	}

	if err := s.storage.IncrementClicks(shortCode); err != nil {
		// логируем, но не возвращаем ошибку
	}

	return url.OriginalURL, nil
}

func (s *URLService) GetStats(shortCode string) (originalURL string, clicks int, createdAt time.Time, err error) {
	url, exists := s.storage.Get(shortCode)
	if !exists {
		return "", 0, time.Time{}, errors.New("short code not found")
	}

	return url.OriginalURL, url.Clicks, url.CreatedAt, nil
}
