include .env

run:
	go run ./cmd/main.go

build:
	go build -o bin/tixgo ./cmd/main.go

create_migration:
	migrate create -ext=sql -dir=migrations/ -seq init_schema

migrate_up:
	migrate -path=migrations/ -database=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable up

migrate_down:
	migrate -path=migrations/ -database=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable down

migrate_force:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make migrate_force VERSION=<version>"; \
		exit 1; \
	fi
	migrate -path=migrations/ -database=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable force $(VERSION)

.PHONY: run build create_migration migrate_up migrate_down migrate_force
