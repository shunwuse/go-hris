package infra

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger structure
type Logger struct {
	*zap.Logger
}

var (
	globalLogger  *Logger
	newLoggerOnce sync.Once
)

// GetLogger returns the global logger instance
func GetLogger() Logger {
	newLoggerOnce.Do(func() {
		logger := newLogger(GetConfig())
		globalLogger = &logger
	})

	return *globalLogger
}

func newLogger(config Config) Logger {
	// get directory path
	dir := filepath.Dir(config.LogOutput)

	// check directory logs exists, if not create it
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	// create logger core based on environment
	var core zapcore.Core
	if config.Environment == "development" {
		core = createDevelopmentCore(config)
	} else {
		core = createProductionCore(config)
	}

	// create logger with caller
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// set global logger
	zap.ReplaceGlobals(logger)

	return Logger{logger}
}

func createDevelopmentCore(config Config) zapcore.Core {
	encoderConfig := createEncoderConfig()
	file := createFileWriter(config.LogOutput)

	// file core with JSON format
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zapcore.DebugLevel,
	)

	// console core with colored output
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	// combine both cores
	return zapcore.NewTee(fileCore, consoleCore)
}

func createProductionCore(config Config) zapcore.Core {
	encoderConfig := createEncoderConfig()
	file := createFileWriter(config.LogOutput)

	// only file output with JSON format
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zapcore.InfoLevel,
	)
}

func createEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func createFileWriter(logOutput string) *os.File {
	file, err := os.OpenFile(logOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// fallback to stderr if file cannot be opened
		return os.Stderr
	}
	return file
}
