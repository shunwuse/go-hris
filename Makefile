# Variables
MIGRATIONS_DIR = ./migrations
VERSION = $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS = -X github.com/shunwuse/go-hris/internal/infra/app.Version=$(VERSION)
DOCKER_IMAGE = go-hris:latest

# Colors for help
BLUE = \033[0;34m
NC = \033[0m

.PHONY: help \
	run run-worker \
	fmt modernize gen \
	test \
	migrate-create \
	docker-build


all: help

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-25s$(NC) %s\n", $$1, $$2}'

## Development targets
run: ## Run the api server locally
	go run -ldflags "$(LDFLAGS)" ./cmd/server

run-worker: ## Run the worker locally
	go run -ldflags "$(LDFLAGS)" ./cmd/worker

fmt: ## Format the code
	go fmt ./...

modernize: ## Apply x/tools modernizations
	go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest -fix ./...

# go install github.com/google/wire/cmd/wire@latest
gen: ## Generate wire dependencies
	wire ./cmd/server
	wire ./cmd/worker
	go mod tidy

## Test targets
test: ## Run unit tests
	go test -v ./...

## Database targets
migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## Docker targets
docker-build: ## Build Docker image
	docker buildx build --platform linux/amd64 --build-arg VERSION=$(VERSION) -t $(DOCKER_IMAGE) --load -f build/package/Dockerfile .
