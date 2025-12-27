package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type RequestLoggerMiddleware struct {
	logger *infra.Logger
}

func NewRequestLoggerMiddleware(logger *infra.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{
		logger: logger,
	}
}

func (m *RequestLoggerMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code and size.
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Process request.
			next.ServeHTTP(ww, r)

			// Log request details.
			latency := time.Since(start)

			m.logger.WithContext(r.Context()).Info("http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()),
				zap.Duration("latency", latency),
				zap.String("latency_human", latency.String()),
			)
		})
	}
}

func (m *RequestLoggerMiddleware) Setup(router chi.Router) {
	router.Use(m.Handler())
}
