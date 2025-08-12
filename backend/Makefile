.PHONY: help build test clean run docker-build docker-run migrate-up migrate-down mod-tidy

# デフォルトタスク
help:
	@echo "Available commands:"
	@echo "  build       - Build the application binary"
	@echo "  test        - Run all tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  run         - Run the application locally"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  migrate-up  - Run database migrations up"
	@echo "  migrate-down- Run database migrations down"
	@echo "  mod-tidy    - Clean up go.mod and go.sum"

# ビルド関連
build:
	go build -o bin/server ./cmd/server

test:
	go test -v ./...

test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# 実行関連
run:
	go run ./cmd/server

# Docker関連
docker-build:
	docker build -t lstep-backend .

docker-run:
	docker run -p 8080:8080 lstep-backend

# データベース関連（goose使用想定）
migrate-up:
	goose -dir migrations postgres "$(DB_DSN)" up

migrate-down:
	goose -dir migrations postgres "$(DB_DSN)" down

migrate-status:
	goose -dir migrations postgres "$(DB_DSN)" status

# 依存関係管理
mod-tidy:
	go mod tidy
	go mod download

# リント関連（golangci-lint使用想定）
lint:
	golangci-lint run

# 開発用（ホットリロード）
dev:
	air

# テストデータベースセットアップ
test-db-up:
	docker run --name lstep-test-db -e POSTGRES_PASSWORD=testpass -e POSTGRES_DB=testdb -p 5433:5432 -d postgres:15

test-db-down:
	docker stop lstep-test-db || true
	docker rm lstep-test-db || true