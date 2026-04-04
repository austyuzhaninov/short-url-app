package app

import (
	"fmt"
	"log"
	"short-url-app/internal/endpoint"
	"short-url-app/internal/pkg/middleware"
	"short-url-app/internal/service"

	"github.com/labstack/echo/v5"
	echomv "github.com/labstack/echo/v5/middleware"
)

type App struct {
	endpoint *endpoint.Endpoint
	service  *service.Service
	echo     *echo.Echo
}

func New() (*App, error) {
	app := &App{}

	app.service = service.New()
	app.endpoint = endpoint.New(app.service)
	app.echo = echo.New()

	// Middlewares
	app.echo.Use(echomv.RequestLogger()) // use the RequestLogger mw with slog logger
	app.echo.Use(echomv.Recover())       // recover panics as errors for proper error handling

	// Custom middlewares
	app.echo.Use(middleware.RoleCheck)

	//Handlers
	app.echo.POST("/shorten", app.endpoint.Shorten) // Генерация короткой ссылки

	return app, nil
}

func (app *App) Run() error {
	fmt.Printf("Server is running")

	err := app.echo.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
