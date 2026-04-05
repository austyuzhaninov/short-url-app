package app

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"short-url-app/internal/endpoint"
	"short-url-app/internal/pkg/config"
	"short-url-app/internal/pkg/validator"
	"short-url-app/internal/service"
	"short-url-app/internal/storage"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomv "github.com/labstack/echo/v4/middleware"
)

type App struct {
	endpoint *endpoint.URLEndpoint
	service  *service.URLService
	storage  storage.Storage
	echo     *echo.Echo
	config   *config.Config
}

func New(cfg *config.Config) (*App, error) {
	// Инициализация storage (снизу вверх)
	store, err := storage.NewMemoryStorage(cfg.StorageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to init storage: %w", err)
	}

	// Инициализация service
	svc := service.New(store, cfg.BaseURL)

	// Инициализация endpoint
	ep := endpoint.New(svc)

	// Инициализация HTTP сервера
	server := echo.New()
	server.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Second
	server.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Second

	// Подключаем кастомный валидатор
	server.Validator = validator.NewEchoValidator()

	// Middlewares
	server.Use(echomv.RequestLogger())
	server.Use(echomv.Recover())
	server.Use(echomv.CORS())

	// Роуты
	server.POST("/shorten", ep.Shorten)
	server.GET("/:code", ep.Redirect)
	server.GET("/stats/:code", ep.GetStats)

	// Health check
	server.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	return &App{
		endpoint: ep,
		service:  svc,
		storage:  store,
		echo:     server,
		config:   cfg,
	}, nil
}

func (a *App) Run() error {
	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запуск сервера в горутине
	errCh := make(chan error, 1)
	go func() {
		log.Printf("Server is running on %s", a.config.Port)
		if err := a.echo.Start(a.config.Port); err != nil {
			errCh <- err
		}
	}()

	// Ожидание сигнала или ошибки
	select {
	case <-ctx.Done():
		// Получен сигнал завершения
		log.Println("Shutting down gracefully...")

		// Контекст с таймаутом для завершения
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Сохраняем данные
		if err := a.storage.SaveToFile(); err != nil {
			log.Printf("Failed to save storage: %v", err)
		}

		// Останавливаем сервер
		if err := a.echo.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error during shutdown: %v", err)
			return err
		}

		log.Println("Server stopped gracefully")
		return nil

	case err := <-errCh:
		// Ошибка при запуске сервера
		log.Printf("Failed to start server: %v", err)
		return err
	}
}
