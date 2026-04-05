package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"short-url-app/internal/endpoint/dto"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// testValidator для тестов
type testValidator struct {
	validator *validator.Validate
}

func (tv *testValidator) Validate(i interface{}) error {
	return tv.validator.Struct(i)
}

// MockURLService реализует интерфейс service.URLServiceInterface
type MockURLService struct {
	ShortenURLFunc     func(originalURL, userID string) (string, string, error)
	GetOriginalURLFunc func(shortCode string) (string, error)
	GetStatsFunc       func(shortCode string) (string, int, time.Time, error)
}

func (m *MockURLService) ShortenURL(originalURL, userID string) (string, string, error) {
	if m.ShortenURLFunc != nil {
		return m.ShortenURLFunc(originalURL, userID)
	}
	return "abc123", "http://localhost:8080/abc123", nil
}

func (m *MockURLService) GetOriginalURL(shortCode string) (string, error) {
	if m.GetOriginalURLFunc != nil {
		return m.GetOriginalURLFunc(shortCode)
	}
	return "", nil
}

func (m *MockURLService) GetStats(shortCode string) (string, int, time.Time, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(shortCode)
	}
	return "", 0, time.Time{}, nil
}

// setupTestServer создаёт тестовый сервер с моком
func setupTestServer(mockService *MockURLService) *echo.Echo {
	e := echo.New()

	// Настраиваем валидатор
	validate := validator.New()
	e.Validator = &testValidator{validator: validate}

	// Создаём endpoint с моком
	ep := New(mockService)

	// Регистрируем роуты
	e.POST("/shorten", ep.Shorten)
	e.GET("/:code", ep.Redirect)
	e.GET("/stats/:code", ep.GetStats)

	return e
}

// ========== TESTS FOR SHORTEN ==========

func TestShorten_Success(t *testing.T) {
	// Arrange
	mockService := &MockURLService{
		ShortenURLFunc: func(url, userID string) (string, string, error) {
			return "xyz789", "http://localhost:8080/xyz789", nil
		},
	}
	e := setupTestServer(mockService)

	reqBody := dto.ShortenRequest{
		URL:    "https://example.com",
		UserID: "testuser",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}

	var resp dto.ShortenResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if resp.ShortCode != "xyz789" {
		t.Errorf("Expected short_code 'xyz789', got '%s'", resp.ShortCode)
	}
	if resp.ShortURL != "http://localhost:8080/xyz789" {
		t.Errorf("Expected short_url 'http://localhost:8080/xyz789', got '%s'", resp.ShortURL)
	}
}

func TestShorten_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := &MockURLService{}
	e := setupTestServer(mockService)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestShorten_EmptyURL(t *testing.T) {
	// Arrange
	mockService := &MockURLService{}
	e := setupTestServer(mockService)

	reqBody := dto.ShortenRequest{
		URL:    "",
		UserID: "testuser",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	// Arrange
	mockService := &MockURLService{}
	e := setupTestServer(mockService)

	reqBody := dto.ShortenRequest{
		URL:    "not-a-valid-url",
		UserID: "testuser",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestShorten_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockURLService{
		ShortenURLFunc: func(url, userID string) (string, string, error) {
			return "", "", echo.ErrInternalServerError
		},
	}
	e := setupTestServer(mockService)

	reqBody := dto.ShortenRequest{
		URL:    "https://example.com",
		UserID: "testuser",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

// ========== TESTS FOR REDIRECT ==========

func TestRedirect_Success(t *testing.T) {
	// Arrange
	mockService := &MockURLService{
		GetOriginalURLFunc: func(shortCode string) (string, error) {
			return "https://example.com", nil
		},
	}
	e := setupTestServer(mockService)

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status 301, got %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("Expected Location 'https://example.com', got '%s'", location)
	}
}

func TestRedirect_EmptyCode(t *testing.T) {
	// Arrange
	mockService := &MockURLService{}
	e := setupTestServer(mockService)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	// Arrange
	mockService := &MockURLService{
		GetOriginalURLFunc: func(shortCode string) (string, error) {
			return "", echo.ErrNotFound
		},
	}
	e := setupTestServer(mockService)

	req := httptest.NewRequest(http.MethodGet, "/notexist", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}

// ========== TESTS FOR STATS ==========

func TestStats_Success(t *testing.T) {
	// Arrange
	mockService := &MockURLService{
		GetStatsFunc: func(shortCode string) (string, int, time.Time, error) {
			createdAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
			return "https://example.com", 42, createdAt, nil
		},
	}
	e := setupTestServer(mockService)

	req := httptest.NewRequest(http.MethodGet, "/stats/abc123", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var resp dto.StatsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if resp.OriginalURL != "https://example.com" {
		t.Errorf("Expected OriginalURL 'https://example.com', got '%s'", resp.OriginalURL)
	}
	if resp.Clicks != 42 {
		t.Errorf("Expected Clicks 42, got %d", resp.Clicks)
	}
	if resp.CreatedAt.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", resp.CreatedAt.Year())
	}
}

func TestStats_EmptyCode(t *testing.T) {
	mockService := &MockURLService{
		GetOriginalURLFunc: func(shortCode string) (string, error) {
			// Этот вызов произойдёт, потому что /stats/ попадёт в Redirect
			return "https://example.com", nil
		},
	}
	e := setupTestServer(mockService)

	req := httptest.NewRequest(http.MethodGet, "/stats/", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	// /stats/ без кода должен обрабатываться как редирект или 404?
	// Сейчас он попадает в Redirect и возвращает 301
	if rec.Code != http.StatusNotFound && rec.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status 404 or 301, got %d", rec.Code)
	}
}

func TestStats_NotFound(t *testing.T) {
	// Arrange
	mockService := &MockURLService{
		GetStatsFunc: func(shortCode string) (string, int, time.Time, error) {
			return "", 0, time.Time{}, echo.ErrNotFound
		},
	}
	e := setupTestServer(mockService)

	req := httptest.NewRequest(http.MethodGet, "/stats/notexist", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}
