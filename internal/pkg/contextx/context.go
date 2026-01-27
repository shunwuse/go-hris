package contextx

import (
	"context"

	"github.com/shunwuse/go-hris/internal/constants"
)

// GetTraceID extracts the TraceID from the context.
func GetTraceID(ctx context.Context) string {
	val, ok := ctx.Value(constants.TraceID).(string)
	if !ok {
		return ""
	}

	return val
}

// WithTraceID returns a new context with the TraceID.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, constants.TraceID, traceID)
}

// GetSpanID extracts the SpanID from the context.
func GetSpanID(ctx context.Context) string {
	val, ok := ctx.Value(constants.SpanID).(string)
	if !ok {
		return ""
	}

	return val
}

// WithSpanID returns a new context with the SpanID.
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, constants.SpanID, spanID)
}
