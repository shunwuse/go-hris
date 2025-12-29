package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/utils"
)

type TraceMiddleware struct {
	logger *infra.Logger
}

func NewTraceMiddleware(
	logger *infra.Logger,
) *TraceMiddleware {
	return &TraceMiddleware{
		logger: logger,
	}
}

func (m *TraceMiddleware) Setup(router chi.Router) {
	router.Use(m.Handler())
}

func (m *TraceMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get("X-Trace-Id")
			if traceID == "" {
				traceID = utils.NewTraceID()
			}

			// Set trace ID in response header.
			w.Header().Set("X-Trace-Id", traceID)

			// Store trace ID in context.
			ctx := context.WithValue(r.Context(), constants.TraceID, traceID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
