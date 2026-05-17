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
	./scripts/run.sh

run-worker: ## Run the worker locally
	./scripts/run-worker.sh

fmt: ## Format the code
	go fmt ./...

modernize: ## Apply x/tools modernizations
	go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest -fix ./...

# go install github.com/google/wire/cmd/wire@latest
gen: ## Generate wire and ent dependencies
	./scripts/gen.sh

## Test targets
test: ## Run unit tests
	go test -v ./...

## Database targets
migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	./scripts/migrate-create.sh "$(name)"

## Docker targets
docker-build: ## Build Docker image
	./scripts/docker-build.sh
