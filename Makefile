# variables
migrations_dir = ./migrations
sqlite_db = ./test.db

server:
	go run ./cmd/server/main.go

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

.PHONY: server \
	migrate-create migrate-up migrate-down \
	go-migrate-up go-migrate-down \
	swagger
