package middlewares

import (
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra"
)

type TraceMiddleware struct {
	logger      *infra.Logger
	entropyPool *sync.Pool
}

func NewTraceMiddleware(
	logger *infra.Logger,
) *TraceMiddleware {
	pool := &sync.Pool{
		New: func() any {
			return ulid.Monotonic(rand.Reader, 0)
		},
	}

	return &TraceMiddleware{
		logger:      logger,
		entropyPool: pool,
	}
}

func (m TraceMiddleware) Setup(router chi.Router) {
	router.Use(m.Handler())
}

func (m TraceMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get("X-Trace-Id")
			if traceID == "" {
				entropy := m.entropyPool.Get().(io.Reader)

				traceID = ulid.MustNew(ulid.Now(), entropy).String()

				// Return entropy to pool for reuse.
				m.entropyPool.Put(entropy)
			}

			// Set trace ID in response header.
			w.Header().Set("X-Trace-Id", traceID)

			// Store trace ID in context.
			ctx := context.WithValue(r.Context(), constants.TraceID, traceID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
