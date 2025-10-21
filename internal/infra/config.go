package infra

import (
	"log"
	"sync"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Environment string `koanf:"env"`
	ServerPort  string `koanf:"server_port"`
	LogOutput   string `koanf:"log_output"`

	SqliteDBPath string `koanf:"sqlite_db_path"`

	JWTSecret string `koanf:"jwt_secret"`
}

var (
	globalConfig   *Config
	loadConfigOnce sync.Once
)

// GetConfig returns a copy of the config
func GetConfig() Config {
	loadConfigOnce.Do(func() {
		globalConfig = loadConfig()
	})
	return *globalConfig
}

func loadConfig() *Config {
	config := &Config{}

	k := koanf.New(".")

	// read env file
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	// unmarshal config
	if err := k.Unmarshal("", config); err != nil {
		log.Fatalf("Unable to unmarshal config: %v", err)
	}

	return config
}
