# variables
migrations_dir = ./migrations
sqlite_db = ./test.db

server:
	go run ./cmd/server/main.go

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

.PHONY: server \
	migrate-create migrate-up migrate-down \
	go-migrate-up go-migrate-down
