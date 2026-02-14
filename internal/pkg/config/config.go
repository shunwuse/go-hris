package config

import (
	"log"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	instance atomic.Pointer[Config]
	loadOnce sync.Once
)

// Load returns the current configuration, loading it if necessary.
func Load() (*Config, error) {
	var err error
	loadOnce.Do(func() {
		var cfg *Config
		cfg, err = load()
		if err == nil {
			instance.Store(cfg)
		}
	})

	if err != nil {
		return nil, err
	}

	return instance.Load(), nil
}

func load() (*Config, error) {
	cfg := &Config{}

	k := koanf.New(".")

	// [STAGE 1] >>> Load 'default.env' (Default Values).
	if err := k.Load(file.Provider("configs/default.env"), dotenv.Parser()); err != nil {
		log.Printf("info: configs/default.env file not found, skipping: %v", err)
	}

	// [STAGE 2] >>> Load '.env' (Local Overrides).
	if err := k.Load(file.Provider("configs/.env"), dotenv.Parser()); err != nil {
		log.Printf("info: configs/.env file not found, skipping: %v", err)
	}

	// [STAGE 3] >>> Load 'Environment Variables' (System Overrides).
	// Prefix "APP_", example: APP_SERVER_PORT will map to server_port.
	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "APP_"))
	}), nil); err != nil {
		log.Printf("warning: failed to load environment variables: %v", err)
	}

	// Unmarshal configuration.
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
