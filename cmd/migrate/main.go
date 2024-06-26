package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("No command provided")
		os.Exit(1)
	}

	m, err := migrate.New("file://./migrations", "sqlite3://./test.db")
	if err != nil {
		log.Fatalf("migrate connection error: %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatalf("migrate up error: %v", err)
		}
	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("migrate down error: %v", err)
		}
	default:
		log.Fatalf("Usage: %s [up|down]", os.Args[0])
	}
}
