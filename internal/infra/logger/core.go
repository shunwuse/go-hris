package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	consoleTimeFormat = "15:04:05.000"
)

// createCore creates the zapcore.Core with unified logic.
func createCore(cfg Config, level zap.AtomicLevel, fields []zap.Field, enableConsole bool) zapcore.Core {
	var cores []zapcore.Core

	// File Core.
	if cfg.FilePath != "" {
		writer := createFileWriter(cfg)
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(createEncoderConfig()),
			zapcore.AddSync(writer),
			level,
		)

		cores = append(cores, fileCore)
	}

	// Console Core.
	if enableConsole {
		// Console core with colored output.
		consoleConfig := zap.NewDevelopmentEncoderConfig()
		consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleConfig.EncodeTime = zapcore.TimeEncoderOfLayout(consoleTimeFormat)

		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)

		cores = append(cores, consoleCore)
	}

	// If no cores defined, default to stdout JSON.
	if len(cores) == 0 {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(createEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// Combine cores and apply global fields.
	return zapcore.NewTee(cores...).With(fields)
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

func createFileWriter(cfg Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,    // megabytes
		MaxBackups: cfg.MaxBackups, // number of backups
		MaxAge:     cfg.MaxAge,     // days
		Compress:   cfg.Compress,   // compress old files
	}
}
