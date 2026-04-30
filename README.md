# TaskFlow API

![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

REST API для управления задачами (To-Do сервис) на Go.

## Стек

- **Go** 1.22+
- **Fiber** — HTTP-фреймворк
- **GORM** — ORM
- **SQLite** — база данных (локально)
- **ozzo-validation** — валидация запросов
- **swaggo** — Swagger UI для go
- [**SwagGen**](https://github.com/Tuxi4k/SwagGen/) — моя утилита для генерации Swagger-документации
- **Taskfile** — управление командами проекта

## Быстрый старт

```bash
# Установка зависимостей
go mod tidy

# Запуск сервера
task dev
```

Сервер поднимается на `http://0.0.0.0:3000`.

## Структура проекта

```
taskflow/
├── cmd/
│   └── main.go                     # Точка входа
├── internal/
│   ├── database/
│   │   └── db.go                   # Подключение к БД, миграции
│   └── modules/
│       └── task/
│           ├── models.go           # Сущность Task, Status
│           ├── dto.go              # Input DTO + валидация
│           ├── repository.go       # Слой доступа к данным
│           ├── service.go          # Бизнес-логика
│           ├── handler.go          # HTTP-хендлеры
│           └── tests/              # Тесты
│               ├── service_test.go
│               └── integration_test.go
├── Taskfile.yml                    # Команды проекта
├── go.mod
└── go.sum
```

## API Endpoints

| Метод  | Endpoint   | Описание                         |
| ------ | ---------- | -------------------------------- |
| GET    | /tasks     | Список задач (фильтр `?status=`) |
| GET    | /tasks/:id | Задача по ID                     |
| POST   | /tasks     | Создать задачу                   |
| PATCH  | /tasks/:id | Частично обновить задачу         |
| DELETE | /tasks/:id | Удалить задачу                   |

## Примеры запросов и ответов

### 1 Создать задачу

**Запрос** `POST /tasks`

```json
{
  "title": "Купить молоко",
  "status": "todo"
}
```

**Ответ** `201 Created`

```json
{
  "id": 1,
  "title": "Купить молоко",
  "status": "todo",
  "created_at": "2026-04-30T05:00:00Z",
  "updated_at": "2026-04-30T05:00:00Z"
}
```

### 2 Получить список задач

**Запрос** `GET /tasks`

**Ответ** `200 OK`

```json
[
  {
    "id": 1,
    "title": "Купить молоко",
    "status": "todo",
    "created_at": "2026-04-30T05:00:00Z",
    "updated_at": "2026-04-30T05:00:00Z"
  },
  {
    "id": 2,
    "title": "Позвонить клиенту",
    "status": "doing",
    "created_at": "2026-04-30T05:10:00Z",
    "updated_at": "2026-04-30T05:12:00Z"
  }
]
```

### 3 Получить задачу по ID

**Запрос** `GET /tasks/1`

**Ответ** `200 OK`

```json
{
  "id": 1,
  "title": "Купить молоко",
  "status": "todo",
  "created_at": "2026-04-30T05:00:00Z",
  "updated_at": "2026-04-30T05:00:00Z"
}
```

### 4 Обновить задачу

**Запрос** `PATCH /tasks/1`

```json
{
  "status": "done"
}
```

**Ответ** `200 OK`

```json
{
  "id": 1,
  "title": "Купить молоко",
  "status": "done",
  "created_at": "2026-04-30T05:00:00Z",
  "updated_at": "2026-04-30T05:20:00Z"
}
```

### 5 Удалить задачу

**Запрос** `DELETE /tasks/1`

**Ответ** `204 No Content`

```json
{}
```

### Пример ошибки

**Ответ** `400 Bad Request` `404 Not Found` `500 Internal Server Error`

```json
{
  "error": "message"
}
```

## Валидация

| Поле   | Правила                                     |
| ------ | ------------------------------------------- |
| title  | обязательно, 3-200 символов                 |
| status | `todo`, `doing`, `done` (по умолчанию todo) |

Ошибки возвращаются в формате:

```json
{
  "title": "обязательное поле",
  "status": "допустимые статусы: todo, doing, done"
}
```

## Swagger

Документация доступна при запуске сервера:

```
http://localhost:3000/swagger/index.html
```

Генерация JSON для продакшена:

```bash
swag init --ot json -g cmd/main.go
```

## Команды

```bash
task test        # Запуск всех тестов
task dev        # Запуск сервера разработки
```

## Тестирование

```bash
# Unit-тесты сервиса
go test ./internal/modules/task/tests/ -v -run TestService

# Интеграционные тесты
go test ./internal/modules/task/tests/ -v -run TestIntegration

# Все тесты
task test
```

## Коды ответов

| Код | Значение                  |
| --- | ------------------------- |
| 200 | Успех                     |
| 201 | Создано                   |
| 204 | Удалено                   |
| 400 | Неверные данные           |
| 404 | Не найдено                |
| 422 | Ошибка валидации          |
| 500 | Внутренняя ошибка сервера |

## Лицензия

[MIT](License)
