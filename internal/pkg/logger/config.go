package logger

import (
	"go.uber.org/zap"
)

// Config defines the configuration for the logger.
type Config struct {
	Level      string `koanf:"log_level"`
	FilePath   string `koanf:"log_file_path"`
	MaxSize    int    `koanf:"log_max_size"`
	MaxBackups int    `koanf:"log_max_backups"`
	MaxAge     int    `koanf:"log_max_age"`
	Compress   bool   `koanf:"log_compress"`
}

type options struct {
	config        Config
	fields        []zap.Field
	enableConsole bool
}

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	return Config{
		Level: "info",
	}
}
