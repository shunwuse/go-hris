package main

import (
	"log"

	"github.com/shunwuse/go-hris/internal/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log := initLogger(cfg)

	worker := InitializeWorker(cfg, log)

	worker.Run()
}
