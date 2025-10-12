# Subscriptions API

REST API для управления подписками пользователей, написанный на Go с использованием Gin фреймворка и PostgreSQL.

## Быстрый запуск

### Предварительные требования

- **Docker** и **Docker Compose**
- **Make** (для удобства использования команд)
- **golang-migrate** (для работы с миграциями БД)

Установка golang-migrate:
```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Windows
choco install golang-migrate
```

### Запуск проекта

1. **Клонируйте репозиторий:**
```bash
git clone <repository-url>
cd subscriptions-api
```

2. **Запустите Docker-контейнеры:**
```bash
make docker-run
```

Эта команда:
- Остановит и удалит существующие контейнеры
- Пересоберет образы без кэша
- Запустит PostgreSQL и API в контейнерах

3. **Примените миграции базы данных:**
```bash
# Дождитесь полного запуска контейнеров (30-60 секунд), затем:
make migrate-up
```

4. **Готово!**
    - API доступен по адресу: http://localhost:8080
    - Swagger документация: http://localhost:8080/swagger/index.html

## Доступные команды

### Docker
```bash
make docker-run    # Запустить проект в Docker
```

### База данных
```bash
make migrate-up      # Применить миграции
make migrate-down    # Откатить миграции
make migrate-version # Показать текущую версию миграций
make migrate-force   # Принудительно установить версию (make migrate-force V=1)
```

### Разработка
```bash
make run               # Запустить локально (требует локальную БД)
make build            # Собрать бинарник
make generate-swagger # Обновить Swagger документацию
```

## Архитектура проекта

```
subscriptions-api/
├── cmd/subscriptions-api/     # Точка входа приложения
├── internal/                  # Приватный код приложения
│   ├── config/                # Конфигурация
│   ├── db/postgres/           # Подключение к БД
│   ├── domain/                # Доменные сущности и ошибки
│   ├── dto/                   # Data Transfer Objects
│   ├── handler/               # HTTP обработчики
│   ├── repository/            # Слой данных
│   ├── router/                # Маршрутизация
│   └── service/               # Бизнес-логика
├── configs/                   # YAML конфигурации
├── migrations/                # SQL миграции
├── docs/                      # Swagger документация
├── docker-compose.yaml        # Docker Compose конфигурация
├── Dockerfile                 # Образ приложения
└── Makefile                   # Команды для разработки
```

## Конфигурация

### Окружения
- **Docker:** `configs/config.docker.yaml`
- **Локальное:** `configs/config.local.yaml`

### Настройки базы данных (Docker)
- **База данных:** subscriptions_db
- **Пользователь:** postgres
- **Пароль:** postgres123
- **Порт:** 5433 (на хосте) → 5432 (в контейнере)

## API Endpoints

После запуска проекта доступна интерактивная Swagger документация по адресу:
**http://localhost:8080/swagger/index.html**

Основные endpoints:
- `GET /api/v1/subscriptions` - Получить все подписки
- `POST /api/v1/subscriptions` - Создать подписку
- `GET /api/v1/subscriptions/{id}` - Получить подписку по ID
- `PUT /api/v1/subscriptions/{id}` - Обновить подписку
- `DELETE /api/v1/subscriptions/{id}` - Удалить подписку
- `GET /api/v1/subscriptions/total` - Получить общую стоимость подписок за период
- `GET /health` - Проверка состояния API

## Заметки для проверки

Это тестовое задание, поэтому:
- Конфиги включены в репозиторий для удобства
- Логи выводятся в debug режиме
- Swagger документация включена для тестирования API