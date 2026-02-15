package middlewares

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/ports/infra"
)

type ExceptionMiddleware struct {
	alerter infra.Alerter
}

func NewExceptionMiddleware(alerter infra.Alerter) *ExceptionMiddleware {
	return &ExceptionMiddleware{
		alerter: alerter,
	}
}

func (m *ExceptionMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			// Alert for 5xx errors
			// Since RecoveryMiddleware is now outside of ExceptionMiddleware,
			// a panic will bypass this logic entirely, avoiding double alerts.
			if ww.Status() >= 500 {
				traceID := contextx.GetTraceID(r.Context())

				_ = m.alerter.Send(r.Context(), infra.Message{
					Level:   infra.LevelError,
					TraceID: traceID,
					Title:   fmt.Sprintf("HTTP %d Error", ww.Status()),
					Content: fmt.Sprintf("Method: %s, Path: %s, Status: %d", r.Method, r.URL.Path, ww.Status()),
				})
			}
		})
	}
}

func (m *ExceptionMiddleware) Setup(router chi.Router) {
	router.Use(m.Handler())
}
