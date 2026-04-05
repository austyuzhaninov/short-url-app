# Переменные
BINARY_NAME=url-shortener
BUILD_DIR=bin
DOCKER_IMAGE=url-shortener
DOCKER_CONTAINER=url-shortener
COMPOSE_FILE=docker-compose.yml

# Цвета для вывода
GREEN=\033[0;32m
RED=\033[0;31m
YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: help
help: ## Показать помощь
	@echo "$(GREEN)Available commands:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

# ========== Локальный запуск ==========

.PHONY: run
run: ## Запустить приложение локально
	@echo "$(GREEN)Running application...$(NC)"
	go run cmd/server/main.go

.PHONY: build
build: ## Собрать бинарник
	@echo "$(GREEN)Building binary...$(NC)"
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go
	@echo "$(GREEN)Binary built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

.PHONY: clean
clean: ## Очистить бинарники и временные файлы
	@echo "$(GREEN)Cleaning...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f storage.json
	rm -f coverage.out
	rm -rf data
	@echo "$(GREEN)Clean completed$(NC)"

# ========== Тесты ==========

.PHONY: test
test: ## Запустить все тесты
	@echo "$(GREEN)Running all tests...$(NC)"
	go test -v ./...

.PHONY: test-unit
test-unit: ## Запустить юнит-тесты
	@echo "$(GREEN)Running unit tests...$(NC)"
	go test -v ./internal/...

.PHONY: test-integration
test-integration: ## Запустить интеграционные тесты
	@echo "$(GREEN)Running integration tests...$(NC)"
	go test -v ./tests/integration/...

.PHONY: test-cover
test-cover: ## Запустить тесты с покрытием
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

.PHONY: test-race
test-race: ## Запустить тесты с детектором гонок
	@echo "$(GREEN)Running tests with race detector...$(NC)"
	CGO_ENABLED=1 go test -race -v ./...

# ========== Docker ==========

.PHONY: docker-build
docker-build: ## Собрать Docker образ
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE) -f docker/Dockerfile .

.PHONY: docker-run
docker-run: ## Запустить Docker контейнер
	@echo "$(GREEN)Running Docker container...$(NC)"
	docker run -d \
		--name $(DOCKER_CONTAINER) \
		-p 8080:8080 \
		-e PORT=8080 \
		-e STORAGE_FILE=/root/data/storage.json \
		-e BASE_URL=http://localhost:8080 \
		-v $(PWD)/data:/root/data \
		$(DOCKER_IMAGE)

.PHONY: docker-stop
docker-stop: ## Остановить Docker контейнер
	@echo "$(GREEN)Stopping Docker container...$(NC)"
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

.PHONY: docker-logs
docker-logs: ## Показать логи Docker контейнера
	docker logs -f $(DOCKER_CONTAINER)

.PHONY: docker-shell
docker-shell: ## Зайти в shell Docker контейнера
	docker exec -it $(DOCKER_CONTAINER) sh

# ========== Docker Compose ==========

.PHONY: compose-up
compose-up: ## Запустить через docker-compose
	@echo "$(GREEN)Starting with docker-compose...$(NC)"
	docker-compose -f $(COMPOSE_FILE) up -d --build

.PHONY: compose-down
compose-down: ## Остановить docker-compose
	@echo "$(GREEN)Stopping docker-compose...$(NC)"
	docker-compose -f $(COMPOSE_FILE) down

.PHONY: compose-logs
compose-logs: ## Показать логи docker-compose
	docker-compose -f $(COMPOSE_FILE) logs -f

.PHONY: compose-restart
compose-restart: compose-down compose-up ## Перезапустить docker-compose

# ========== Разработка ==========

.PHONY: fmt
fmt: ## Форматировать код
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...

.PHONY: lint
lint: ## Запустить линтер (требуется golangci-lint)
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run

.PHONY: tidy
tidy: ## Обновить зависимости
	@echo "$(GREEN)Tidying dependencies...$(NC)"
	go mod tidy

.PHONY: deps
deps: tidy ## Установить зависимости
	@echo "$(GREEN)Installing dependencies...$(NC)"
	go mod download

.PHONY: dev
dev: ## Запустить в режиме разработки (требуется air)
	@echo "$(GREEN)Running in development mode...$(NC)"
	air

# ========== Инициализация ==========

.PHONY: init
init: ## Инициализировать проект (создать .env, data dir)
	@echo "$(GREEN)Initializing project...$(NC)"
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(GREEN)Created .env file from .env.example$(NC)"; \
	fi
	@mkdir -p data
	@mkdir -p $(BUILD_DIR)
	@echo "$(GREEN)Project initialized$(NC)"

# ========== Утилиты ==========

.PHONY: version
version: ## Показать версии инструментов
	@echo "$(YELLOW)Go version:$(NC) $(shell go version)"
	@echo "$(YELLOW)Docker version:$(NC) $(shell docker --version)"
	@echo "$(YELLOW)Docker Compose version:$(NC) $(shell docker-compose --version)"

.PHONY: info
info: ## Показать информацию о проекте
	@echo "$(GREEN)Project: URL Shortener$(NC)"
	@echo "$(YELLOW)Binary name:$(NC) $(BINARY_NAME)"
	@echo "$(YELLOW)Build dir:$(NC) $(BUILD_DIR)"
	@echo "$(YELLOW)Docker image:$(NC) $(DOCKER_IMAGE)"
	@echo "$(YELLOW)Docker container:$(NC) $(DOCKER_CONTAINER)"