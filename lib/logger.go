package lib

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger structure
type Logger struct {
	*logrus.Logger
}

var globalLogger *Logger

func GetLogger() Logger {
	if globalLogger == nil {
		logger := newLogger()
		globalLogger = &logger
	}

	return *globalLogger
}

func newLogger() Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}

	return Logger{logger}
}
