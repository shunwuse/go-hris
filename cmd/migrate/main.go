package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/shunwuse/go-hris/internal/pkg/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("No command provided")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.SSLMode,
	)

	m, err := migrate.New("file://./migrations", dsn)
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
