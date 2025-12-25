# Variables
APP_NAME = myapp
MIGRATIONS_DIR = ./migrations
DB_URL = ./test.db
VERSION = $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS = -X github.com/shunwuse/go-hris/internal/infra.Version=$(VERSION)

# Colors for help
BLUE = \033[0;34m
NC = \033[0m

.PHONY: help \
	run \
	gen \
	test test-coverage \
	test-integration test-integration-quick test-integration-endpoints \
	build build-static \
	migrate-create migrate-up migrate-down \
	go-migrate-up go-migrate-down \
	docker-build docker-run


all: help

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-25s$(NC) %s\n", $$1, $$2}'

## Development targets
run: ## Run the application locally
	go run -ldflags "$(LDFLAGS)" ./cmd/server

# go install github.com/google/wire/cmd/wire@latest
gen: ## Generate wire dependencies
	wire ./cmd/server
	go mod tidy

## Test targets
test: ## Run unit tests
	go test -v ./...

test-coverage: ## Run tests with coverage report
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-integration: ## Run all integration tests
	./scripts/run_tests_with_server.sh

test-integration-quick: ## Run quick integration tests
	./scripts/quick_test.sh

test-integration-endpoints: ## Run endpoint integration tests
	./scripts/test_endpoints.sh

build: ## Build the application
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) ./cmd/server

build-static: ## Build the application statically
	go build -ldflags "$(LDFLAGS) -extldflags '-static'" -o $(APP_NAME) ./cmd/server

## Database targets
migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

migrate-up: ## Run migrations up using migrate CLI
	migrate -path $(MIGRATIONS_DIR) -database "sqlite3://$(DB_URL)" up

migrate-down: ## Run migrations down using migrate CLI
	migrate -path $(MIGRATIONS_DIR) -database "sqlite3://$(DB_URL)" down

go-migrate-up: ## Run migrations up using go run
	go run ./cmd/migrate/main.go up

go-migrate-down: ## Run migrations down using go run
	go run ./cmd/migrate/main.go down

## Docker targets
docker-build: ## Build Docker image
	docker buildx build --platform linux/amd64 --build-arg VERSION=$(VERSION) -t go-hris:latest .

docker-run: ## Run Docker container
	docker run --rm -p 8080:8080 go-hris:latest
