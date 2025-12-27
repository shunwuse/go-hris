package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type RecoveryMiddleware struct {
	logger *infra.Logger
}

func NewRecoveryMiddleware(logger *infra.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger: logger,
	}
}

func (m *RecoveryMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if err == http.ErrAbortHandler {
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(err)
					}

					// Log the panic with stack trace.
					m.logger.WithContext(r.Context()).Error("panic recovered",
						zap.Any("error", err),
						zap.String("stack", string(debug.Stack())),
					)

					// Set connection to close to prevent connection reuse after panic.
					if r.Header.Get("Connection") != "Upgrade" {
						w.Header().Set("Connection", "close")
					}

					// Return structured error response.
					response.Error(w, errors.ErrInternalError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func (m *RecoveryMiddleware) Setup(router chi.Router) {
	router.Use(m.Handler())
}
