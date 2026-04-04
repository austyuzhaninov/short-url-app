package main

import (
	"log"
	"short-url-app/internal/app"
	"short-url-app/internal/pkg/config"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Создание приложения
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Запуск приложения
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
