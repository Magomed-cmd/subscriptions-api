include .env
export

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
MIGRATE_CMD=migrate -path migrations -database "$(DB_URL)"

run:
	go run cmd/subscriptions-api/main.go

build:
	go build -o subscriptions-api cmd/subscriptions-api/main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

migrate-force:
	$(MIGRATE_CMD) force $(or $(V),)

migrate-version:
	$(MIGRATE_CMD) version

docker-run:
	docker compose -f docker-compose.yml down -v
	docker compose -f docker-compose.yml build --no-cache
	docker compose -f docker-compose.yml up

generate-swagger:
	swag init -g cmd/subscriptions-api/main.go -o docs/swagger
