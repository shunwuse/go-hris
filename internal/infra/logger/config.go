package logger

import (
	"go.uber.org/zap"
)

// Config defines the configuration for the logger.
type Config struct {
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
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
