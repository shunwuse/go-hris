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

	server := InitializeServer(cfg)

	server.Run()
}
