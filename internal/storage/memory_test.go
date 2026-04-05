package storage

import (
	"os"
	"short-url-app/internal/models/entity"
	"testing"
)

// TestSaveAndGet проверяет базовое сохранение и получение
func TestSaveAndGet(t *testing.T) {
	// Создаём временный файл
	tmpFile, err := os.CreateTemp("", "storage_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	storage, err := NewMemoryStorage(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	url := entity.URL{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		UserID:      "user1",
		Clicks:      0,
	}

	// Сохраняем
	err = storage.Save(url)
	if err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// Получаем
	saved, exists := storage.Get("abc123")
	if !exists {
		t.Error("URL not found after save")
	}
	if saved.OriginalURL != url.OriginalURL {
		t.Errorf("Expected %s, got %s", url.OriginalURL, saved.OriginalURL)
	}
}

// TestIncrementClicks проверяет увеличение счётчика
func TestIncrementClicks(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "storage_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	storage, err := NewMemoryStorage(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	url := entity.URL{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		UserID:      "user1",
		Clicks:      0,
	}

	storage.Save(url)

	// Инкрементируем
	err = storage.IncrementClicks("abc123")
	if err != nil {
		t.Errorf("IncrementClicks failed: %v", err)
	}

	// Проверяем
	saved, _ := storage.Get("abc123")
	if saved.Clicks != 1 {
		t.Errorf("Expected clicks=1, got %d", saved.Clicks)
	}
}

// TestExists проверяет проверку существования кода
func TestExists(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "storage_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	storage, err := NewMemoryStorage(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	url := entity.URL{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		UserID:      "user1",
	}

	storage.Save(url)

	if !storage.Exists("abc123") {
		t.Error("Exists should return true for existing code")
	}
	if storage.Exists("notexist") {
		t.Error("Exists should return false for non-existing code")
	}
}

// TestPersistence проверяет сохранение данных в файл и загрузку
func TestPersistence(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "storage_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Создаём первое хранилище и сохраняем данные
	storage1, err := NewMemoryStorage(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	url := entity.URL{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		UserID:      "user1",
		Clicks:      5,
	}

	err = storage1.Save(url)
	if err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// Создаём второе хранилище (как после перезапуска)
	storage2, err := NewMemoryStorage(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем, что данные загрузились
	saved, exists := storage2.Get("abc123")
	if !exists {
		t.Error("Data not persisted after restart")
	}
	if saved.OriginalURL != "https://example.com" {
		t.Errorf("Expected original URL 'https://example.com', got '%s'", saved.OriginalURL)
	}
	if saved.Clicks != 5 {
		t.Errorf("Expected clicks=5, got %d", saved.Clicks)
	}
	if saved.UserID != "user1" {
		t.Errorf("Expected user_id='user1', got '%s'", saved.UserID)
	}
}

// TestGetAll проверяет получение всех записей
func TestGetAll(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "storage_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	storage, err := NewMemoryStorage(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Сохраняем несколько записей
	urls := []entity.URL{
		{ShortCode: "aaa", OriginalURL: "https://site1.com", UserID: "user1"},
		{ShortCode: "bbb", OriginalURL: "https://site2.com", UserID: "user2"},
		{ShortCode: "ccc", OriginalURL: "https://site3.com", UserID: "user3"},
	}

	for _, url := range urls {
		storage.Save(url)
	}

	// Получаем все
	all := storage.GetAll()
	if len(all) != 3 {
		t.Errorf("Expected 3 items, got %d", len(all))
	}

	// Проверяем каждый
	for _, expected := range urls {
		saved, exists := all[expected.ShortCode]
		if !exists {
			t.Errorf("ShortCode '%s' not found", expected.ShortCode)
			continue
		}
		if saved.OriginalURL != expected.OriginalURL {
			t.Errorf("For %s: expected URL '%s', got '%s'",
				expected.ShortCode, expected.OriginalURL, saved.OriginalURL)
		}
	}
}
