package infra

import (
	"log/slog"
	"os"
	"path/filepath"
)

// Logger structure
type Logger struct {
}

var globalLogger *Logger

func GetLogger() Logger {
	if globalLogger == nil {
		logger := newLogger(NewEnv())
		globalLogger = &logger
	}

	return *globalLogger
}

func newLogger(env Env) Logger {
	// get directory path
	dir := filepath.Dir(env.LogOutput)

	// check directory logs exists, if not create it
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	file, err := os.OpenFile(env.LogOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Warn("Failed to open log file, using stderr", "error", err)
		file = os.Stderr
	}

	// create JSON handler with debug level
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	})

	// set as the global default logger
	slog.SetDefault(slog.New(handler))

	return Logger{}
}
