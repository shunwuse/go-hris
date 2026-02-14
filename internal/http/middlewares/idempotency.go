package middlewares

import (
	"bytes"
	stderrors "errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shunwuse/go-hris/internal/infra/config"
	"github.com/shunwuse/go-hris/internal/infra/idempotency"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"go.uber.org/zap"
)

type IdempotencyMiddleware struct {
	config  *config.Config
	manager *idempotency.Manager
	logger  *logger.Logger
}

func NewIdempotencyMiddleware(cfg *config.Config, manager *idempotency.Manager, logger *logger.Logger) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		config:  cfg,
		manager: manager,
		logger:  logger,
	}
}

func (m *IdempotencyMiddleware) Setup(router chi.Router) {
	router.Use(m.Handler())
}

// Handler returns a middleware with the default timeout from config.
func (m *IdempotencyMiddleware) Handler() func(http.Handler) http.Handler {
	return m.HandlerWithTTL(time.Duration(m.config.Service.IdempotencyExpireMinutes) * time.Minute)
}

// HandlerWithTTL returns a middleware with a custom timeout.
func (m *IdempotencyMiddleware) HandlerWithTTL(ttl time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Only apply to state-changing operations.
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			// Get Idempotency Key (Reuse Trace ID from context).
			key := contextx.GetTraceID(ctx)
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check cache.
			record, err := m.manager.Get(ctx, key)
			if err == nil && record != nil {
				m.logger.Info("idempotency cache hit",
					zap.String("key", key),
				)
				for k, v := range record.Header {
					w.Header().Set(k, v)
				}
				w.Header().Set("X-Idempotency-Hit", "true")
				w.WriteHeader(record.Status)
				w.Write(record.Body) //nolint:errcheck
				return
			}

			if err != nil && !stderrors.Is(err, idempotency.ErrNotFound) {
				m.logger.Error("failed to get idempotency record", zap.Error(err))
			}

			// Proxy and Capture.
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			bodyBuf := &bytes.Buffer{}
			ww.Tee(bodyBuf) // Synchronously copy response body to bodyBuf

			next.ServeHTTP(ww, r)

			// Store result (only if status is 2xx or 4xx, avoid caching 5xx transient errors).
			if ww.Status() >= 200 && ww.Status() < 500 {
				header := make(map[string]string)
				if ct := ww.Header().Get("Content-Type"); ct != "" {
					header["Content-Type"] = ct
				}

				if err := m.manager.Set(ctx, key, &idempotency.Record{
					Status: ww.Status(),
					Body:   bodyBuf.Bytes(),
					Header: header,
				}, ttl); err != nil {
					m.logger.Error("failed to set idempotency record", zap.Error(err))
				}
			}
		})
	}
}
