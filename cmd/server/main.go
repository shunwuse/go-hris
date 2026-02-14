package main

import (
	"log"

	"github.com/shunwuse/go-hris/internal/infra/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log := initLogger(cfg)

	server := InitializeServer(cfg, log)

	server.Run()
}
