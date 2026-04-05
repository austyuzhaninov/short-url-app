package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"short-url-app/internal/app"
	"short-url-app/internal/endpoint/dto"
	"short-url-app/internal/pkg/config"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestApp создаёт полноценное приложение для интеграционных тестов
func setupTestApp(t *testing.T) (*app.App, string) {
	// Создаём временный файл для storage
	tmpFile, err := os.CreateTemp("", "integration_test_*.json")
	require.NoError(t, err)
	tmpFilePath := tmpFile.Name()
	tmpFile.Close()

	// Настраиваем конфиг для тестов
	cfg := &config.Config{
		Port:         ":0", // случайный порт
		StorageFile:  tmpFilePath,
		BaseURL:      "http://test.local",
		ReadTimeout:  30,
		WriteTimeout: 30,
	}

	// Создаём приложение
	application, err := app.New(cfg)
	require.NoError(t, err)

	return application, tmpFilePath
}

// setupTestAppWithConfig создаёт приложение из готового конфига
func setupTestAppWithConfig(t *testing.T, cfg *config.Config) *app.App {
	application, err := app.New(cfg)
	require.NoError(t, err)
	return application
}

// cleanup удаляет временные файлы
func cleanup(t *testing.T, filePath string) {
	os.Remove(filePath)
}

// TestIntegration_CreateAndRedirect полный сценарий: создание → редирект → статистика
func TestIntegration_CreateAndRedirect(t *testing.T) {
	// Arrange
	application, storagePath := setupTestApp(t)
	defer cleanup(t, storagePath)
	e := application.Echo

	// Act: Создаём короткую ссылку
	createReq := dto.ShortenRequest{
		URL:    "https://golang.org",
		UserID: "testuser",
	}
	bodyBytes, _ := json.Marshal(createReq)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert: Проверяем создание
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createResp dto.ShortenResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createResp)
	require.NoError(t, err)
	assert.NotEmpty(t, createResp.ShortCode)
	assert.Equal(t, "http://test.local/"+createResp.ShortCode, createResp.ShortURL)

	shortCode := createResp.ShortCode

	// Act: Переходим по короткой ссылке
	req = httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert: Проверяем редирект
	assert.Equal(t, http.StatusMovedPermanently, rec.Code)
	assert.Equal(t, "https://golang.org", rec.Header().Get("Location"))

	// Act: Получаем статистику
	req = httptest.NewRequest(http.MethodGet, "/stats/"+shortCode, nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert: Проверяем статистику
	assert.Equal(t, http.StatusOK, rec.Code)

	var statsResp dto.StatsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &statsResp)
	require.NoError(t, err)
	assert.Equal(t, "https://golang.org", statsResp.OriginalURL)
	assert.Equal(t, 1, statsResp.Clicks) // один переход
	assert.False(t, statsResp.CreatedAt.IsZero())
}

// TestIntegration_InvalidURL возвращает ошибку на невалидный URL
func TestIntegration_InvalidURL(t *testing.T) {
	// Arrange
	application, storagePath := setupTestApp(t)
	defer cleanup(t, storagePath)
	e := application.Echo

	createReq := dto.ShortenRequest{
		URL:    "not-a-valid-url",
		UserID: "testuser",
	}
	bodyBytes, _ := json.Marshal(createReq)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestIntegration_NotFound возвращает 404 для несуществующего кода
func TestIntegration_NotFound(t *testing.T) {
	// Arrange
	application, storagePath := setupTestApp(t)
	defer cleanup(t, storagePath)
	e := application.Echo

	// Act: Переход по несуществующему коду
	req := httptest.NewRequest(http.MethodGet, "/notexist", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Act: Статистика по несуществующему коду
	req = httptest.NewRequest(http.MethodGet, "/stats/notexist", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestIntegration_MultipleRequests создаёт несколько ссылок
func TestIntegration_MultipleRequests(t *testing.T) {
	// Arrange
	application, storagePath := setupTestApp(t)
	defer cleanup(t, storagePath)
	e := application.Echo

	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
	}

	codes := make([]string, 0, len(urls))

	// Act: Создаём несколько ссылок
	for _, url := range urls {
		createReq := dto.ShortenRequest{
			URL:    url,
			UserID: "testuser",
		}
		bodyBytes, _ := json.Marshal(createReq)

		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp dto.ShortenResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		codes = append(codes, resp.ShortCode)
	}

	// Проверяем, что все коды уникальны
	uniqueCodes := make(map[string]bool)
	for _, code := range codes {
		assert.False(t, uniqueCodes[code], "Duplicate code: %s", code)
		uniqueCodes[code] = true
	}
	assert.Len(t, uniqueCodes, len(urls))
}

// TestIntegration_Persistence проверяет сохранение данных между перезапусками
func TestIntegration_Persistence(t *testing.T) {
	// Arrange: Первый запуск
	app1, storagePath := setupTestApp(t)
	defer cleanup(t, storagePath)
	e1 := app1.Echo

	// Создаём ссылку
	createReq := dto.ShortenRequest{
		URL:    "https://persistence-test.com",
		UserID: "testuser",
	}
	bodyBytes, _ := json.Marshal(createReq)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e1.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var createResp dto.ShortenResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createResp)
	require.NoError(t, err)
	shortCode := createResp.ShortCode

	// Делаем один переход
	req = httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
	rec = httptest.NewRecorder()
	e1.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMovedPermanently, rec.Code)

	// Act: "Перезапускаем" приложение — создаём новое приложение с тем же storage файлом
	cfg := &config.Config{
		Port:         ":0",
		StorageFile:  storagePath,
		BaseURL:      "http://test.local",
		ReadTimeout:  30,
		WriteTimeout: 30,
	}
	app2 := setupTestAppWithConfig(t, cfg)
	e2 := app2.Echo

	// Assert: Проверяем, что данные сохранились
	req = httptest.NewRequest(http.MethodGet, "/stats/"+shortCode, nil)
	rec = httptest.NewRecorder()
	e2.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var statsResp dto.StatsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &statsResp)
	require.NoError(t, err)
	assert.Equal(t, "https://persistence-test.com", statsResp.OriginalURL)
	assert.Equal(t, 1, statsResp.Clicks)
}

// TestIntegration_ConcurrentRequests проверяет конкурентные запросы
func TestIntegration_ConcurrentRequests(t *testing.T) {
	// Arrange
	application, storagePath := setupTestApp(t)
	defer cleanup(t, storagePath)
	e := application.Echo

	n := 50
	successCount := 0

	// Act: Отправляем конкурентные запросы на создание
	for i := 0; i < n; i++ {
		go func(i int) {
			createReq := dto.ShortenRequest{
				URL:    "https://example.com",
				UserID: "user",
			}
			bodyBytes, _ := json.Marshal(createReq)

			req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code == http.StatusCreated {
				successCount++
			}
		}(i)
	}

	// Ждём завершения всех горутин
	time.Sleep(2 * time.Second)

	// Assert: Проверяем, что все запросы успешны
	assert.Equal(t, n, successCount)
}
