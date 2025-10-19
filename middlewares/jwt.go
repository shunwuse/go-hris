package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/ports/service"
)

type JWTMiddleware struct {
	logger      lib.Logger
	authService service.AuthService
}

func NewJWTMiddleware(
	logger lib.Logger,
	authService service.AuthService,
) JWTMiddleware {
	return JWTMiddleware{
		logger:      logger,
		authService: authService,
	}
}

func (m JWTMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{
					"error": "Authorization header required",
				})
				return
			}

			texts := strings.Split(authHeader, " ")
			if len(texts) != 2 {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{
					"error": "Invalid Authorization header",
				})
				return
			}

			if texts[0] != "Bearer" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{
					"error": "Bearer token required",
				})
				return
			}

			token := texts[1]

			claims, err := m.authService.AuthenticateToken(r.Context(), token)
			if err != nil {
				m.logger.Errorf("authenticating token failed: %v", err)
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{
					"error": "Invalid token",
				})
				return
			}

			ctx := context.WithValue(r.Context(), constants.JWTClaims, claims.TokenPayload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
