package endpoint

import (
	"net/http"
	"short-url-app/internal/endpoint/dto"
	"short-url-app/internal/service"

	"github.com/labstack/echo/v4"
)

type URLEndpoint struct {
	service *service.URLService
}

func New(service *service.URLService) *URLEndpoint {
	return &URLEndpoint{
		service: service,
	}
}

func (e *URLEndpoint) Shorten(c echo.Context) error {
	var req dto.ShortenRequest

	// Bind + валидация в одной строке!
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Автоматическая валидация по тегам
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	shortCode, shortURL, err := e.service.ShortenURL(req.URL, req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.ShortenResponse{
		ShortCode: shortCode,
		ShortURL:  shortURL,
	})
}

func (e *URLEndpoint) Redirect(c echo.Context) error {
	shortCode := c.Param("code")
	if shortCode == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "code is required"})
	}

	originalURL, err := e.service.GetOriginalURL(shortCode)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (e *URLEndpoint) GetStats(c echo.Context) error {
	shortCode := c.Param("code")
	if shortCode == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "code is required"})
	}

	originalURL, clicks, createdAt, err := e.service.GetStats(shortCode)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.StatsResponse{
		OriginalURL: originalURL,
		Clicks:      clicks,
		CreatedAt:   createdAt,
	})
}
