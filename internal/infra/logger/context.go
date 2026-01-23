package logger

import (
	"context"

	"github.com/shunwuse/go-hris/internal/constants"
	"go.uber.org/zap"
)

// WithContext returns a logger with trace ID from context if available.
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	if traceID, ok := ctx.Value(constants.TraceID).(string); ok {
		return l.With(zap.String("trace_id", traceID))
	}

	return l.Logger
}
