package endpoint

import (
	"net/http"
	"short-url-app/internal/endpoint/dto"
	"short-url-app/internal/pkg/response"
	"short-url-app/internal/service"

	"github.com/labstack/echo/v4"
)

type URLEndpoint struct {
	service service.URLServiceInterface
}

func New(service service.URLServiceInterface) *URLEndpoint {
	return &URLEndpoint{
		service: service,
	}
}

func (e *URLEndpoint) Shorten(c echo.Context) error {
	var req dto.ShortenRequest

	if err := response.BindAndValidate(c, &req); err != nil {
		return err
	}

	shortCode, shortURL, err := e.service.ShortenURL(req.URL, req.UserID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, http.StatusCreated, dto.ShortenResponse{
		ShortCode: shortCode,
		ShortURL:  shortURL,
	})
}

func (e *URLEndpoint) Redirect(c echo.Context) error {
	shortCode, err := response.GetParam(c, "code")
	if err != nil {
		return err
	}

	originalURL, err := e.service.GetOriginalURL(shortCode)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (e *URLEndpoint) GetStats(c echo.Context) error {
	shortCode, err := response.GetParam(c, "code")
	if err != nil {
		return err
	}

	originalURL, clicks, createdAt, err := e.service.GetStats(shortCode)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, http.StatusOK, dto.StatsResponse{
		OriginalURL: originalURL,
		Clicks:      clicks,
		CreatedAt:   createdAt,
	})
}
