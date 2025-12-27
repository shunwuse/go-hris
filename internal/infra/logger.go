package infra

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/shunwuse/go-hris/internal/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	consoleTimeFormat = "15:04:05.000"
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

	// Define common fields.
	fields := []zapcore.Field{
		zap.String("instance_id", InstanceID),
		zap.String("hostname", Hostname),
		zap.Int("pid", os.Getpid()),
		zap.String("environment", config.Environment),
	}

	// Create logger core based on environment.
	var core zapcore.Core
	if config.Environment == constants.EnvDevelopment {
		core = createDevelopmentCore(config, fields)
	} else {
		core = createProductionCore(config, fields)
	}

	// Create logger with caller.
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Set global logger.
	zap.ReplaceGlobals(logger)

	return Logger{logger}
}

func createDevelopmentCore(config Config, fields []zapcore.Field) zapcore.Core {
	encoderConfig := createEncoderConfig()
	writer := createFileWriter(config)

	// File core with JSON format and metadata.
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		zapcore.DebugLevel,
	).With(fields)

	// Console core with colored output.
	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleConfig.EncodeTime = zapcore.TimeEncoderOfLayout(consoleTimeFormat)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	// Combine both cores.
	return zapcore.NewTee(fileCore, consoleCore)
}

func createProductionCore(config Config, fields []zapcore.Field) zapcore.Core {
	encoderConfig := createEncoderConfig()
	writer := createFileWriter(config)

	// Only file output with JSON format and metadata.
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		zapcore.InfoLevel,
	).With(fields)
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

func createFileWriter(config Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   config.LogOutput,
		MaxSize:    config.LogMaxSize,    // megabytes
		MaxBackups: config.LogMaxBackups, // number of backups
		MaxAge:     config.LogMaxAge,     // days
		Compress:   config.LogCompress,   // compress old files
	}
}
