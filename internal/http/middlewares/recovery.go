package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra/alerter"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"go.uber.org/zap"
)

type RecoveryMiddleware struct {
	logger  *logger.Logger
	alerter alerter.Alerter
}

func NewRecoveryMiddleware(log *logger.Logger, alerter alerter.Alerter) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger:  log,
		alerter: alerter,
	}
}

func (m *RecoveryMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if err == http.ErrAbortHandler {
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged.
						panic(err)
					}

					// Build stack trace.
					stack := string(debug.Stack())

					// Log the panic.
					m.logger.WithContext(r.Context()).Error("panic recovered",
						zap.Any("error", err),
						zap.String("stack", stack),
					)

					// Send Critical Alert.
					traceID, _ := r.Context().Value(constants.TraceID).(string)

					_ = m.alerter.Send(r.Context(), alerter.Message{
						Level:      alerter.LevelCritical,
						TraceID:    traceID,
						Title:      "Panic Recovered",
						Content:    fmt.Sprintf("Method: %s, Path: %s, Error: %v", r.Method, r.URL.Path, err),
						StackTrace: stack,
					})

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
