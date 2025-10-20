package infra

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Logger structure
type Logger struct {
	*logrus.Logger
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
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	// get directory path
	dir := filepath.Dir(env.LogOutput)

	// check directory logs exists, if not create it
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	file, err := os.OpenFile(env.LogOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}

	return Logger{logger}
}
