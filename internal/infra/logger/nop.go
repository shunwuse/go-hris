package logger

import (
	"go.uber.org/zap"
)

// NewNopLogger returns a no-operation logger for testing purposes.
func NewNopLogger() *Logger {
	return &Logger{
		Logger: zap.NewNop(),
		level:  zap.NewAtomicLevel(),
	}
}
