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

type Config struct {
	Environment   string `koanf:"env"`
	ServerPort    string `koanf:"server_port"`
	LogOutput     string `koanf:"log_output"`
	LogMaxSize    int    `koanf:"log_max_size"`    // megabytes
	LogMaxBackups int    `koanf:"log_max_backups"` // number of backups
	LogMaxAge     int    `koanf:"log_max_age"`     // days
	LogCompress   bool   `koanf:"log_compress"`    // compress old files

	SqliteDBPath string `koanf:"sqlite_db_path"`

	JWTSecret               string `koanf:"jwt_secret"`
	JWTExpireMinutes        int    `koanf:"jwt_expire_minutes"`
	JWTRefreshExpireMinutes int    `koanf:"jwt_refresh_expire_minutes"`

	UseMiniredis  bool   `koanf:"use_miniredis"`
	RedisAddr     string `koanf:"redis_addr"`
	RedisPassword string `koanf:"redis_password"`
	RedisDB       int    `koanf:"redis_db"`
}

var (
	instance atomic.Pointer[Config]
	loadOnce sync.Once
)

// Get returns the current configuration.
func Get() Config {
	loadOnce.Do(func() {
		instance.Store(load())
	})

	return *instance.Load()
}

func load() *Config {
	config := &Config{}

	k := koanf.New(".")

	// [STAGE 1] >>> Load 'default.env' (Default Values)
	if err := k.Load(file.Provider("configs/default.env"), dotenv.Parser()); err != nil {
		log.Printf("info: configs/default.env file not found, skipping: %v", err)
	}

	// [STAGE 2] >>> Load '.env' (Local Overrides)
	if err := k.Load(file.Provider("configs/.env"), dotenv.Parser()); err != nil {
		log.Printf("info: configs/.env file not found, skipping: %v", err)
	}

	// [STAGE 3] >>> Load 'Environment Variables' (System Overrides)
	// Prefix "APP_", example: APP_SERVER_PORT will map to server_port.
	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "APP_"))
	}), nil); err != nil {
		log.Printf("warning: failed to load environment variables: %v", err)
	}

	// Unmarshal configuration.
	if err := k.Unmarshal("", config); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return config
}
