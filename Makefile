DB_URL=postgres://postgres:postgres123@localhost:5433/subscriptions_db?sslmode=disable
MIGRATE_CMD=migrate -path migrations -database "$(DB_URL)"

run:
	go run cmd/subscriptions-api/main.go

build:
	go build -o subscriptions-api cmd/subscriptions-api/main.go

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
	docker-compose down -v
	docker-compose build --no-cache
	docker-compose up

generate-swagger:
	swag init -g cmd/subscriptions-api/main.go