# Subscriptions API

REST API для управления подписками пользователей, написанный на Go с использованием Gin, PostgreSQL и принципов Clean Architecture.

## Описание проекта

Сервис предоставляет полный CRUD для управления платными подписками: создание, получение, обновление, удаление, фильтрация по пользователю/сервису и подсчёт суммарной стоимости за период.

## Архитектура

Проект следует **многоуровневой архитектуре** с чётким разделением ответственности:

```text
┌──────────────────────────────┐
│        Handler Layer         │
│  (HTTP, validation, routing) │
└──────────┬───────────────────┘
           │
┌──────────▼───────────────────┐
│        Service Layer         │
│    (Business logic, rules)   │
└──────────┬───────────────────┘
           │
┌──────────▼───────────────────┐
│      Repository Layer        │
│  (Database, query building)  │
└──────────┬───────────────────┘
           │
┌──────────▼───────────────────┐
│         PostgreSQL           │
└──────────────────────────────┘
```

## Структура проекта

```text
subscriptions-api/
├── cmd/subscriptions-api/     # Точка входа, graceful shutdown
├── internal/
│   ├── config/                # Конфигурация (Viper, YAML)
│   ├── db/postgres/           # Подключение к БД (pgxpool)
│   ├── domain/
│   │   ├── entity/            # Доменные сущности
│   │   └── errors/            # Доменные ошибки + маппинг PG-кодов
│   ├── dto/                   # Request/Response объекты, конвертация
│   ├── handler/               # HTTP-хендлеры (Gin), интерфейс сервиса
│   ├── repository/            # SQL-запросы (squirrel query builder)
│   ├── router/                # Регистрация маршрутов
│   └── service/               # Бизнес-логика, валидация
├── configs/                   # YAML конфигурации (local/docker)
├── docs/swagger/              # Swagger/OpenAPI сгенерированные файлы
├── migrations/                # SQL миграции (golang-migrate)
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── .env.example
```

## Быстрый старт

### Требования

- Docker и Docker Compose
- `golang-migrate` для миграций

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

### Установка и запуск

1. **Клонируйте репозиторий**

   ```bash
   git clone <repository-url>
   cd subscriptions-api
   ```

2. **Настройте переменные окружения**

   ```bash
   cp .env.example .env
   ```

3. **Запустите Docker-контейнеры**

   ```bash
   make docker-run
   ```

4. **Примените миграции**

   ```bash
   make migrate-up
   ```

5. **Готово**

   - API: [http://localhost:8080](http://localhost:8080)
   - Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## API Endpoints

| Метод | Путь | Описание |
| ----- | ---- | -------- |
| `POST` | `/api/v1/subscriptions` | Создать подписку |
| `GET` | `/api/v1/subscriptions` | Список подписок (с фильтрами) |
| `GET` | `/api/v1/subscriptions/:id` | Получить подписку по ID |
| `PUT` | `/api/v1/subscriptions/:id` | Обновить подписку |
| `DELETE` | `/api/v1/subscriptions/:id` | Удалить подписку |
| `GET` | `/api/v1/subscriptions/total` | Суммарная стоимость за период |
| `GET` | `/health` | Статус сервиса |

### Параметры фильтрации (GET /subscriptions)

| Параметр | Тип | Описание |
| -------- | --- | -------- |
| `user_id` | UUID | Фильтр по пользователю |
| `service_name` | string | Фильтр по названию сервиса |
| `start_date` | `2006-01-02` | Дата начала диапазона |
| `end_date` | `2006-01-02` | Дата конца диапазона |

### Параметры подсчёта стоимости (GET /subscriptions/total)

| Параметр | Тип | Обязательный | Описание |
| -------- | --- | ------------ | -------- |
| `period_start` | `2006-01-02T15:04:05Z` | Да | Начало периода |
| `period_end` | `2006-01-02T15:04:05Z` | Да | Конец периода |
| `user_id` | UUID | Нет | Фильтр по пользователю |
| `service_name` | string | Нет | Фильтр по сервису |

## Доступные команды

```bash
# Запуск
make run              # Локально (требует локальную БД)
make docker-run       # В Docker

# Сборка
make build            # Собрать бинарник

# Качество кода
make fmt              # go fmt
make vet              # go vet

# Миграции
make migrate-up       # Применить миграции
make migrate-down     # Откатить миграции
make migrate-version  # Текущая версия
make migrate-force V=1  # Принудительно установить версию

# Документация
make generate-swagger # Обновить Swagger docs
```

## Конфигурация

Конфиг выбирается автоматически по окружению:

- **Docker**: `configs/config.docker.yaml`
- **Локально**: `configs/config.local.yaml`

Переменные окружения переопределяют значения из YAML-файлов с префиксом `APP_`.

## Stack технологий

- **Язык**: Go 1.25
- **Framework**: Gin
- **БД**: PostgreSQL 16
- **Драйвер БД**: pgx v5 + pgxpool
- **Query builder**: squirrel
- **Конфигурация**: Viper
- **Логирование**: Uber Zap
- **Миграции**: golang-migrate
- **API Docs**: Swagger/OpenAPI (swaggo)
- **Контейнеризация**: Docker, Docker Compose
