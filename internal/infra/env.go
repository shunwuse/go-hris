package infra

import (
	"log"
	"sync"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Env struct {
	Environment string `koanf:"env"`
	ServerPort  string `koanf:"server_port"`
	LogOutput   string `koanf:"log_output"`

	SqliteDBPath string `koanf:"sqlite_db_path"`

	JWTSecret string `koanf:"jwt_secret"`
}

var (
	globalEnv   *Env
	loadEnvOnce sync.Once
)

// GetEnv returns a copy of the config
func GetEnv() Env {
	loadEnvOnce.Do(func() {
		globalEnv = loadEnv()
	})
	return *globalEnv
}

func loadEnv() *Env {
	env := &Env{}

	k := koanf.New(".")

	// read env file
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	// unmarshal env
	if err := k.Unmarshal("", env); err != nil {
		log.Fatalf("Unable to unmarshal env: %v", err)
	}

	return env
}
