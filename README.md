
# Структура проекта

```text
text
url-shortener/
├── cmd/
│   └── server/
│       └── main.go                 # Точка входа
├── internal/
│   ├── handlers/
│   │   └── url.go                  # HTTP handlers
│   ├── storage/
│   │   ├── memory.go               # In-memory хранилище
│   │   └── file.go                 # Работа с файлом
│   ├── models/
│   │   └── url.go                  # Структуры данных
│   ├── config/
│   │   └── config.go               # Конфигурация
│   └── generator/
│       └── code.go                 # Генератор коротких кодов
├── pkg/
│   └── logger/
│       └── logger.go               # Логирование (опционально)
├── tests/
│   └── integration_test.go         # Интеграционные тесты
├── storage.json                    # Файл для хранения данных (создаётся при запуске)
├── go.mod
├── go.sum
└── README.md
```

