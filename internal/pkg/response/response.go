package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Success возвращает успешный JSON ответ
func Success(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, data)
}

// Error возвращает JSON ошибку
func Error(c echo.Context, status int, errMsg string) error {
	return c.JSON(status, map[string]string{"error": errMsg})
}

// BadRequest возвращает 400 ошибку
func BadRequest(c echo.Context, errMsg string) error {
	return Error(c, http.StatusBadRequest, errMsg)
}

// NotFound возвращает 404 ошибку
func NotFound(c echo.Context, errMsg string) error {
	return Error(c, http.StatusNotFound, errMsg)
}

// BindAndValidate связывает и валидирует запрос
func BindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return BadRequest(c, "invalid request body")
	}
	if err := c.Validate(req); err != nil {
		return BadRequest(c, err.Error())
	}
	return nil
}

// GetParam возвращает параметр из URL с проверкой на пустоту
func GetParam(c echo.Context, paramName string) (string, error) {
	value := c.Param(paramName)
	if value == "" {
		return "", BadRequest(c, paramName+" is required")
	}
	return value, nil
}
