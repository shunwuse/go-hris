package logger

import (
	"go.uber.org/zap"
)

// Option defines a functional option for the logger options.
type Option func(*options)

// WithConfig sets the initial configuration.
func WithConfig(cfg Config) Option {
	return func(o *options) {
		o.config = cfg
	}
}

// WithConsole enables or disables console output.
func WithConsole(enable bool) Option {
	return func(o *options) {
		o.enableConsole = enable
	}
}

// WithFields adds multiple common fields to the logger.
func WithFields(fields ...zap.Field) Option {
	return func(o *options) {
		o.fields = append(o.fields, fields...)
	}
}
