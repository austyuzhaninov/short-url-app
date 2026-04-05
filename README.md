# URL Shortener Service

Сервис для сокращения ссылок с поддержкой статистики переходов.

## Возможности

- Сокращение длинных URL в короткие коды (6 символов)
- Редирект по короткому коду на оригинальный URL
- Статистика переходов по каждой ссылке
- In-memory хранение с синхронизацией через мьютексы
- Персистентность данных в JSON файл
- Graceful shutdown с сохранением данных
- Конфигурация через переменные окружения

## Установка и запуск

### Требования
* Go 1.21 или выше

### Локальный запуск
```bash
# Клонирование репозитория
git clone https://github.com/yourusername/url-shortener.git
cd url-shortener

# Установка зависимостей
go mod tidy

# Запуск сервера
go run cmd/server/main.go
```

### Конфигурация через переменные окружения

| Переменная | Описание | Значение по умолчанию |
|-------|----------|------|
| PORT | Порт сервера | `8080` |
| STORAGE_FILE |Путь к файлу хранения | `./storage.json` |
| BASE_URL | Базовый URL для коротких ссылок | `http://localhost:8080` |
| READ_TIMEOUT | Таймаут чтения (сек) | `30` |
| WRITE_TIMEOUT | Таймаут записи (сек) | `30` |

### Пример запуска конфигурации
```bash
export PORT=9090
export BASE_URL=https://my.domain.com
go run cmd/server/main.go
```

## Тестирование
```bash
# Все тесты с проверкой гонок
go test -race -v ./...

# Только модульные тесты
go test -v ./internal/storage/...
go test -v ./internal/pkg/generator/...

# Интеграционные тесты
go test -v ./tests/integration/...
```

## Примеры использования
```bash
# 1. Создаём короткую ссылку
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://golang.org/doc", "user_id": "alice"}'

# 2. Переходим по ней
curl -v http://localhost:8080/aB3xY9

# 3. Смотрим статистику
curl http://localhost:8080/stats/aB3xY9
```

## Структура проекта
```text
url-shortener/
├── cmd/server/main.go          # Точка входа
├── internal/
│ ├── app/app.go                # Сборка и запуск приложения
│ ├── endpoint/
│   ├── dto/
│   │   └── url.go              # ShortenRequest, ShortenResponse, StatsResponse
│   └── url.go                  # Handlers (использует dto)
│ ├── service/service.go        # Бизнес-логика
│ ├── storage/
│ │ ├── storage.go              # Интерфейс Storage
│ │ └── memory.go               # In-memory реализация
│ ├── models/
│   └── entity/
│       └── url.go              # URL (для storage)
│ └── pkg/
│ ├── config/config.go          # Конфигурация
│ ├── generator/code.go         # Генерация коротких кодов
│ └── validator/echo.go         # Валидация URL
├── tests/                      # Тесты
├── go.mod
├── go.sum
└── storage.json                # Файл с данными
```
## API

### Создание короткой ссылки

```bash
POST /shorten
Content-Type: application/json

{
    "url": "https://example.com",
    "user_id": "alice"
}
```

### Переход по короткой ссылке

```bash
GET /{short_code}
Ответ: HTTP 301 Moved Permanently с Location на оригинальный URL
```

### Статистика
```bash
GET /stats/{short_code}
Ответ (200 OK):
{
    "original_url": "https://example.com",
    "clicks": 42,
    "created_at": "2024-01-15T10:00:00Z"
}
```

### Health check
```bash
GET /health
Ответ (200 OK):
{
    "status": "ok"
}
```

## Архитектура

### Слои приложения


| Слой | Файл | Ответственность |
|-------|----------|------|
| HTTP | `endpoint/url.go` | Парсинг запросов, формирование ответов |
| Бизнес-логика |`service/url.go` | Валидация, генерация кодов, бизнес-правила |
| Хранилище | `storage/` | Интерфейс и реализация хранения |
| Модели | `models/url.go` | Структуры данных (DTO) |
| Утилиты | `pkg/` | Переиспользуемые компоненты |

### Поток данных

```bash
HTTP Request → endpoint → service → storage (интерфейс) → memory → JSON файл
```

### Ключевые принципы

* Инверсия зависимостей: service зависит от интерфейса Storage, а не от конкретной реализации
* Dependency Injection: Все зависимости передаются через конструкторы
* Конкурентность: sync.RWMutex для безопасного доступа к данным
* Graceful shutdown: Корректное завершение с сохранением данных

## Обработка ошибок


| Ситуация | HTTP | статус	Ответ |
|-------|----------|------|
| Невалидный URL | 400 | `{"error": "invalid URL..."}` |
| Несуществующий код |404 | `{"error": "short code not found"}` |
| Некорректный JSON | 400 | `{"error": "invalid request body"}` |
| Отсутствует поле url | 400 | `{"error": "url is required"}` |

