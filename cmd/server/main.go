package main

import (
	"log"

	"github.com/shunwuse/go-hris/internal/pkg/config"
)

func main() {
	mgr := config.NewManager[Config](
		config.WithDotEnv("configs/default.env"),
		config.WithDotEnv("configs/.env"),
		config.WithEnv("APP_", "."),
	)

	if err := mgr.Load(); err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	cfg := mgr.Config()

	log := initLogger(mgr)

	server := InitializeServer(cfg, log)

	server.Run()
}
