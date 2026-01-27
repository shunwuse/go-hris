package logger

import (
	"context"

	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"go.uber.org/zap"
)

// WithContext returns a logger with trace and span IDs from context if available.
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	fields := make([]zap.Field, 0, 2)

	if traceID := contextx.GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if spanID := contextx.GetSpanID(ctx); spanID != "" {
		fields = append(fields, zap.String("span_id", spanID))
	}

	if len(fields) > 0 {
		return l.With(fields...)
	}

	return l.Logger
}
