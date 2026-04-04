package endpoint

import (
	"net/http"
	"short-url-app/internal/models"
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
	var req models.ShortenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.URL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "url is required"})
	}

	resp, err := e.service.ShortenURL(req.URL, req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
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

	stats, err := e.service.GetStats(shortCode)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stats)
}
