package logger

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/infra/config"
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

	level zap.AtomicLevel
}

var (
	instance *Logger
	initOnce sync.Once
)

// L returns the global logger instance.
func L() *Logger {
	initOnce.Do(func() {
		cfg, _ := config.Load()
		logger := newLogger(cfg)
		instance = &logger

		// Register reload hook.
		config.OnChange(func(cfg *config.Config) {
			instance.UpdateLevel(cfg.Log.Level)
		})
	})

	return instance
}

// UpdateLevel dynamically changes the logger level.
func (l *Logger) UpdateLevel(levelStr string) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		l.Error("failed to update log level", zap.String("level", levelStr), zap.Error(err))
		return
	}

	l.level.SetLevel(level)

	l.Info("log level updated", zap.String("new_level", level.String()))
}

func newLogger(cfg *config.Config) Logger {
	// Parse level.
	atomicLevel := zap.NewAtomicLevel()
	if err := atomicLevel.UnmarshalText([]byte(cfg.Log.Level)); err != nil {
		atomicLevel.SetLevel(zapcore.InfoLevel) // Default if invalid.
	}

	// Get directory path.
	dir := filepath.Dir(cfg.Log.FilePath)

	// Check if logs directory exists, if not create it.
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(dir, os.ModePerm); mkErr != nil {
			_, _ = os.Stderr.WriteString("failed to create log directory: " + mkErr.Error() + "\n")
		}
	}

	// Define common fields.
	fields := []zapcore.Field{
		zap.String("instance_id", app.InstanceID),
		zap.String("hostname", app.Hostname),
		zap.Int("pid", os.Getpid()),
		zap.String("environment", cfg.Service.Environment),
	}

	// Create logger core based on environment.
	var core zapcore.Core
	if cfg.Service.Environment == constants.EnvDevelopment {
		core = createDevelopmentCore(cfg, atomicLevel, fields)
	} else {
		core = createProductionCore(cfg, atomicLevel, fields)
	}

	// Create logger with caller.
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Set global logger.
	zap.ReplaceGlobals(logger)

	return Logger{logger, atomicLevel}
}

func createDevelopmentCore(cfg *config.Config, level zap.AtomicLevel, fields []zapcore.Field) zapcore.Core {
	encoderConfig := createEncoderConfig()
	writer := createFileWriter(cfg)

	// File core with JSON format and metadata.
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		level,
	).With(fields)

	// Console core with colored output.
	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleConfig.EncodeTime = zapcore.TimeEncoderOfLayout(consoleTimeFormat)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Combine both cores.
	return zapcore.NewTee(fileCore, consoleCore)
}

func createProductionCore(cfg *config.Config, level zap.AtomicLevel, fields []zapcore.Field) zapcore.Core {
	encoderConfig := createEncoderConfig()
	writer := createFileWriter(cfg)

	// Only file output with JSON format and metadata.
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		level,
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

func createFileWriter(cfg *config.Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,    // megabytes
		MaxBackups: cfg.Log.MaxBackups, // number of backups
		MaxAge:     cfg.Log.MaxAge,     // days
		Compress:   cfg.Log.Compress,   // compress old files
	}
}
