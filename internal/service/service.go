package service

import (
	"errors"
	"fmt"
	"short-url-app/internal/models/dto"
	"short-url-app/internal/models/entity"
	"short-url-app/internal/pkg/generator"
	"short-url-app/internal/pkg/validator"
	"short-url-app/internal/storage"
	"time"
)

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

func (s *URLService) ShortenURL(originalURL, userID string) (*dto.ShortenResponse, error) {
	// Валидация URL
	if !validator.IsValidURL(originalURL) {
		return nil, errors.New("invalid URL: must be valid http:// or https:// URL")
	}

	// Генерация уникального короткого кода
	var shortCode string
	for {
		shortCode = generator.GenerateShortCode()
		if !s.storage.Exists(shortCode) {
			break
		}
	}

	// Создание записи
	url := entity.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		UserID:      userID,
		Clicks:      0,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.storage.Save(url); err != nil {
		return nil, fmt.Errorf("failed to save URL: %w", err)
	}

	return &dto.ShortenResponse{
		ShortCode: shortCode,
		ShortURL:  fmt.Sprintf("%s/%s", s.baseURL, shortCode),
	}, nil
}

func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	url, exists := s.storage.Get(shortCode)
	if !exists {
		return "", errors.New("short code not found")
	}

	// Инкрементируем счётчик переходов
	if err := s.storage.IncrementClicks(shortCode); err != nil {
		// Логируем ошибку, но не возвращаем, чтобы редирект всё равно работал
		// В реальном проекте здесь должен быть логгер
		fmt.Printf("failed to increment clicks: %v\n", err)
	}

	return url.OriginalURL, nil
}

func (s *URLService) GetStats(shortCode string) (*dto.StatsResponse, error) {
	url, exists := s.storage.Get(shortCode)
	if !exists {
		return nil, errors.New("short code not found")
	}

	return &dto.StatsResponse{
		OriginalURL: url.OriginalURL,
		Clicks:      url.Clicks,
		CreatedAt:   url.CreatedAt,
	}, nil
}
