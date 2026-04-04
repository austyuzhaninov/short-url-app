package endpoint

import (
	"short-url-app/internal/service"

	"github.com/labstack/echo/v5"
)

// Service Для дальнейшего тестирования используем интерфейс
type Service interface {
	Shorten() int64
	GetShortCode() int64
	GetStatus() int64
}

type Endpoint struct {
	service Service
}

func New(s *service.Service) *Endpoint {
	return &Endpoint{
		service: s,
	}
}

// Принимаем POST с url, возвращаем ответ с короткой ссылкой
func (e *Endpoint) Shorten(ctx *echo.Context) error {
	//days := e.service.DaysLeft()
	//answer := fmt.Sprintf("Days left: %d", days)
	//
	//err := ctx.String(http.StatusOK, answer)
	//if err != nil {
	//	return err
	//}
	//
	//return nil
}

// Придоставить короткую ссылку
func (e *Endpoint) GetShortCode(ctx *echo.Context) error {

}

func (e *Endpoint) GetStatus(ctx *echo.Context) error {

}
