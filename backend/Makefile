.PHONY: run dev build migrate-up migrate-down

DB_URL ?= postgres://postgres:postgres@localhost:5432/plan_balance?sslmode=disable

run:
	go run cmd/api/main.go

dev:
	# Assuming you have air installed for hot reload
	air

build:
	go build -o bin/api cmd/api/main.go

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down

test:
	go test ./...
