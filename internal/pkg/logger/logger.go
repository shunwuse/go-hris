package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger structure.
type Logger struct {
	*zap.Logger

	level zap.AtomicLevel
}

var (
	instance atomic.Pointer[Logger]
)

// L returns the global logger instance.
func L() *Logger {
	l := instance.Load()
	if l == nil {
		_, _ = fmt.Fprintln(os.Stderr, "CRITICAL: logger.L() called before logger.New(). Application might crash due to nil dereference.")
	}

	return l
}

// New creates a new logger instance.
func New(opts ...Option) *Logger {
	o := &options{
		config:        DefaultConfig(),
		enableConsole: true,
	}

	for _, opt := range opts {
		opt(o)
	}

	cfg := o.config

	// Parse level.
	atomicLevel := zap.NewAtomicLevel()
	if err := atomicLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		atomicLevel.SetLevel(zapcore.InfoLevel)
	}

	// Ensure log directory exists if FilePath is provided.
	if cfg.FilePath != "" {
		dir := filepath.Dir(cfg.FilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if mkErr := os.MkdirAll(dir, os.ModePerm); mkErr != nil {
				_, _ = os.Stderr.WriteString("failed to create log directory: " + mkErr.Error() + "\n")
			}
		}
	}

	// Create logger core.
	core := createCore(cfg, atomicLevel, o.fields, o.enableConsole)

	// Create logger with caller.
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Set global logger.
	zap.ReplaceGlobals(logger)

	log := &Logger{logger, atomicLevel}
	instance.Store(log)

	return log
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
