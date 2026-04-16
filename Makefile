.PHONY: run build test lint migrate-up migrate-down migrate-create tidy

# ── Dev ────────────────────────────────────────────────────
run:
	air

run-no-reload:
	go run ./cmd/server

build:
	go build -o ./bin/server ./cmd/server

# ── Test ───────────────────────────────────────────────────
test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# ── Code Quality ───────────────────────────────────────────
lint:
	go vet ./...

tidy:
	go mod tidy

# ── Database Migration (golang-migrate) ────────────────────
migrate-up:
	migrate -path ./migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" up

migrate-down:
	migrate -path ./migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" down 1

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

# ── Docker (local MySQL) ───────────────────────────────────
docker-db:
	docker-compose up -d mysql

docker-db-stop:
	docker-compose down
