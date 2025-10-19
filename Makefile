# variables
migrations_dir = ./migrations
sqlite_db = ./test.db

server:
	go run ./cmd/server

# go install github.com/google/wire/cmd/wire@latest
wire:
	wire ./cmd/server

# go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(migrations_dir) -seq $$name

migrate-up:
	migrate -path ./migrations -database "sqlite3://$(sqlite_db)" up


migrate-down:
	migrate -path ./migrations -database "sqlite3://$(sqlite_db)" down

go-migrate-up:
	go run ./cmd/migrate/main.go up

go-migrate-down:
	go run ./cmd/migrate/main.go down

# go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	swag init -g ./cmd/server/main.go -o ./docs/swagger

docker-build:
	docker buildx build --platform linux/amd64 -t go-hris:latest .

docker-run:
	docker run --rm -p 8080:8080 go-hris:latest

# Integration Testing
test-integration:
	./scripts/run_tests_with_server.sh

test-integration-quick:
	./scripts/quick_test.sh

test-integration-endpoints:
	./scripts/test_endpoints.sh

# Unit Testing
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: server \
	wire \
	migrate-create migrate-up migrate-down \
	go-migrate-up go-migrate-down \
	swagger \
	docker-build docker-run \
	test test-coverage \
	test-integration test-integration-quick test-integration-endpoints
