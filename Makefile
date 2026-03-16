APP_NAME=payments-service

.PHONY: run test fmt build db-up db-down db-logs db-psql db-schema

run:
	go run ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

build:
	go build -o bin/$(APP_NAME) ./cmd/api

db-up:
	docker compose -f deployments/docker-compose.yml up -d

db-down:
	docker compose -f deployments/docker-compose.yml down

db-logs:
	docker compose -f deployments/docker-compose.yml logs -f postgres

db-psql:
	psql "postgres://payments_service:payments_service@localhost:5432/payments_service?sslmode=disable"

db-schema:
	psql "postgres://payments_service:payments_service@localhost:5432/payments_service?sslmode=disable" -f db/schema.sql
