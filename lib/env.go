package lib

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	Environment string `mapstructure:"ENV"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
	LogOutput   string `mapstructure:"LOG_OUTPUT"`
}

func NewEnv() Env {
	env := Env{}

	// env file path
	viper.SetConfigFile(".env")

	// read env file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading env file, %s", err)
	}

	// unmarshal env
	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return env
}
