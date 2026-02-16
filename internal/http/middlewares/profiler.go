package middlewares

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/infra/app"
)

type ProfilerMiddleware struct {
	token string
}

func NewProfilerMiddleware(cfg *app.ServiceConfig) *ProfilerMiddleware {
	return &ProfilerMiddleware{
		token: cfg.ProfilerToken,
	}
}

func (m *ProfilerMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if m.token != "" {
				token := r.Header.Get("X-Profiler-Token")

				if token != m.token {
					http.Error(w, "Unauthorized: Invalid Profiler Token", http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
