package infra

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/shunwuse/go-hris/internal/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger structure.
type Logger struct {
	*zap.Logger
}

// WithContext returns a logger with trace ID from context if available.
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	if traceID, ok := ctx.Value(constants.TraceID).(string); ok {
		return l.With(zap.String("trace_id", traceID))
	}

	return l.Logger
}

var (
	globalLogger  *Logger
	newLoggerOnce sync.Once
)

// GetLogger returns the global logger instance.
func GetLogger() *Logger {
	newLoggerOnce.Do(func() {
		logger := newLogger(GetConfig())
		globalLogger = &logger
	})

	return globalLogger
}

func newLogger(config Config) Logger {
	// Get directory path.
	dir := filepath.Dir(config.LogOutput)

	// Check if logs directory exists, if not create it.
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(dir, os.ModePerm); mkErr != nil {
			_, _ = os.Stderr.WriteString("failed to create log directory: " + mkErr.Error() + "\n")
		}
	}

	// Create logger core based on environment.
	var core zapcore.Core
	if config.Environment == "development" {
		core = createDevelopmentCore(config)
	} else {
		core = createProductionCore(config)
	}

	// Create logger with caller.
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Set global logger.
	zap.ReplaceGlobals(logger)

	return Logger{logger}
}

func createDevelopmentCore(config Config) zapcore.Core {
	encoderConfig := createEncoderConfig()
	file := createFileWriter(config.LogOutput)

	// File core with JSON format.
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zapcore.DebugLevel,
	)

	// Console core with colored output.
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	// Combine both cores.
	return zapcore.NewTee(fileCore, consoleCore)
}

func createProductionCore(config Config) zapcore.Core {
	encoderConfig := createEncoderConfig()
	file := createFileWriter(config.LogOutput)

	// Only file output with JSON format.
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
		// Fallback to stderr if file cannot be opened.
		return os.Stderr
	}
	return file
}
