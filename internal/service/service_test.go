package service

import (
	"short-url-app/internal/models/entity"
	"testing"
)

// MockStorage — мок для тестирования service
type MockStorage struct {
	urls map[string]entity.URL
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		urls: make(map[string]entity.URL),
	}
}

func (m *MockStorage) Save(url entity.URL) error {
	m.urls[url.ShortCode] = url
	return nil
}

func (m *MockStorage) Get(shortCode string) (entity.URL, bool) {
	url, exists := m.urls[shortCode]
	return url, exists
}

func (m *MockStorage) IncrementClicks(shortCode string) error {
	url, exists := m.urls[shortCode]
	if !exists {
		return nil
	}
	url.Clicks++
	m.urls[shortCode] = url
	return nil
}

func (m *MockStorage) Exists(shortCode string) bool {
	_, exists := m.urls[shortCode]
	return exists
}

func (m *MockStorage) SaveToFile() error {
	return nil
}

func (m *MockStorage) GetAll() map[string]entity.URL {
	return m.urls
}

// TestShortenURL_Valid проверяет создание корректной ссылки
func TestShortenURL_Valid(t *testing.T) {
	mockStorage := NewMockStorage()
	service := New(mockStorage, "http://localhost:8080")

	shortCode, shortURL, err := service.ShortenURL("https://golang.org", "user123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if shortCode == "" {
		t.Error("Expected shortCode not empty")
	}
	if shortURL == "" {
		t.Error("Expected shortURL not empty")
	}

	// Проверяем, что сохранилось в storage
	saved, exists := mockStorage.Get(shortCode)
	if !exists {
		t.Error("URL not saved in storage")
	}
	if saved.OriginalURL != "https://golang.org" {
		t.Errorf("Expected https://golang.org, got %s", saved.OriginalURL)
	}
}

// TestGetOriginalURL_Exists проверяет получение существующей ссылки
func TestGetOriginalURL_Exists(t *testing.T) {
	mockStorage := NewMockStorage()
	service := New(mockStorage, "http://localhost:8080")

	// Сначала создаём ссылку
	shortCode, _, _ := service.ShortenURL("https://example.com", "user123")

	// Потом получаем
	originalURL, err := service.GetOriginalURL(shortCode)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if originalURL != "https://example.com" {
		t.Errorf("Expected https://example.com, got %s", originalURL)
	}
}

// TestGetOriginalURL_NotExists проверяет ошибку на несуществующий код
func TestGetOriginalURL_NotExists(t *testing.T) {
	mockStorage := NewMockStorage()
	service := New(mockStorage, "http://localhost:8080")

	_, err := service.GetOriginalURL("notexist")

	if err == nil {
		t.Error("Expected error for non-existent code, got nil")
	}
}

// TestGetStats проверяет получение статистики
func TestGetStats(t *testing.T) {
	mockStorage := NewMockStorage()
	service := New(mockStorage, "http://localhost:8080")

	// Создаём ссылку
	shortCode, _, _ := service.ShortenURL("https://example.com", "user123")

	// Переходим по ней (увеличиваем clicks)
	service.GetOriginalURL(shortCode)

	// Получаем статистику
	originalURL, clicks, createdAt, err := service.GetStats(shortCode)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if originalURL != "https://example.com" {
		t.Errorf("Expected https://example.com, got %s", originalURL)
	}
	if clicks != 1 {
		t.Errorf("Expected clicks=1, got %d", clicks)
	}
	if createdAt.IsZero() {
		t.Error("Expected createdAt not zero")
	}
}

// TestShortenURL_UniqueCode проверяет генерацию уникальных кодов
func TestShortenURL_UniqueCode(t *testing.T) {
	mockStorage := NewMockStorage()
	service := New(mockStorage, "http://localhost:8080")

	// Создаём несколько ссылок
	codes := make(map[string]bool)
	for i := 0; i < 10; i++ {
		shortCode, _, err := service.ShortenURL("https://example.com", "user123")
		if err != nil {
			t.Errorf("Failed to create URL: %v", err)
		}
		if codes[shortCode] {
			t.Errorf("Duplicate code generated: %s", shortCode)
		}
		codes[shortCode] = true
	}
}

// TestIncrementClicks_OnRedirect проверяет увеличение счётчика при редиректе
func TestIncrementClicks_OnRedirect(t *testing.T) {
	mockStorage := NewMockStorage()
	service := New(mockStorage, "http://localhost:8080")

	// Создаём ссылку
	shortCode, _, _ := service.ShortenURL("https://example.com", "user123")

	// Первый переход
	service.GetOriginalURL(shortCode)
	// Второй переход
	service.GetOriginalURL(shortCode)

	// Проверяем статистику
	_, clicks, _, err := service.GetStats(shortCode)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if clicks != 2 {
		t.Errorf("Expected clicks=2, got %d", clicks)
	}
}
