APP_NAME=todo-api
MAIN=cmd/server/main.go

MIGRATE=migrate -path ./migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable"

run:
	go run $(MAIN)

build:
	go build -o $(APP_NAME) $(MAIN)

tidy:
	go mod tidy

lint:
	golangci-lint run

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-force:
	$(MIGRATE) force

test:
	go test ./...

docs:
	swag init -g cmd/server/main.go -o docs/

.PHONY: run build tidy lint docker-up docker-down migrate-up migrate-down test docs
