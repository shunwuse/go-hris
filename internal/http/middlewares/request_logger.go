package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"go.uber.org/zap"
)

type RequestLoggerMiddleware struct {
	logger *logger.Logger
}

func NewRequestLoggerMiddleware(log *logger.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{
		logger: log,
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
			status := ww.Status()

			// Prepare log fields.
			fields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Int("status", status),
				zap.Int("size", ww.BytesWritten()),
				zap.Duration("latency", latency),
				zap.String("latency_human", latency.String()),
			}

			// Determine log level based on status code.
			log := m.logger.WithContext(r.Context())
			if status >= 500 {
				log.Error("http_request_error", fields...)
			} else if status >= 400 {
				log.Warn("http_request_warn", fields...)
			} else {
				log.Debug("http_request_success", fields...)
			}
		})
	}
}

func (m *RequestLoggerMiddleware) Setup(router chi.Router) {
	router.Use(m.Handler())
}
