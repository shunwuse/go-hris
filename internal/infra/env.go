package infra

import (
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type SqliteConfig struct {
	Database string `koanf:"DB_PATH"`
}

type Env struct {
	Environment string `koanf:"ENV"`
	ServerPort  string `koanf:"SERVER_PORT"`
	LogOutput   string `koanf:"LOG_OUTPUT"`

	Sqlite SqliteConfig `koanf:"SQLITE"`

	JWTSecret string `koanf:"JWT_SECRET"`
}

func NewEnv() Env {
	env := Env{}

	k := koanf.New(".")

	// read env file
	err := k.Load(file.Provider(".env"), dotenv.Parser())
	if err != nil {
		log.Fatalf("Error reading env file, %s", err)
	}

	// convert flat keys with dot(.) to nested structure
	for _, key := range k.Keys() {
		if strings.Contains(key, ".") {
			v := k.String(key)
			k.Delete(key)
			k.Set(key, v)
		}
	}

	// unmarshal env
	err = k.Unmarshal("", &env)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return env
}
