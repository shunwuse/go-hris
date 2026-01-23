package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/utils"
)

type TraceMiddleware struct {
	logger *logger.Logger
}

func NewTraceMiddleware(
	log *logger.Logger,
) *TraceMiddleware {
	return &TraceMiddleware{
		logger: log,
	}
}

func (m *TraceMiddleware) Setup(router chi.Router) {
	router.Use(m.Handler())
}

func (m *TraceMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get or generate Trace ID.
			traceID := r.Header.Get("X-Trace-Id")
			if traceID == "" {
				traceID = utils.NewTraceID()
			}

			// Generate a new Span ID.
			spanID := utils.NewSpanID()

			// Set IDs in response header.
			w.Header().Set("X-Trace-Id", traceID)
			w.Header().Set("X-Span-Id", spanID)

			// Store IDs in context.
			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.SpanID, spanID)
			ctx = context.WithValue(ctx, constants.TraceID, traceID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
