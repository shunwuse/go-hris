package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/services"
)

type JWTMiddleware struct {
	logger      lib.Logger
	authService services.AuthService
}

func NewJWTMiddleware() JWTMiddleware {
	logger := lib.GetLogger()

	// Initialize services
	authService := services.NewAuthService()

	return JWTMiddleware{
		logger:      logger,
		authService: authService,
	}
}

func (m JWTMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(401, gin.H{
				"error": "Authorization header required",
			})
			ctx.Abort()
			return
		}

		texts := strings.Split(authHeader, " ")
		if len(texts) != 2 {
			ctx.JSON(401, gin.H{
				"error": "Invalid Authorization header",
			})
			ctx.Abort()
			return
		}

		if texts[0] != "Bearer" {
			ctx.JSON(401, gin.H{
				"error": "Bearer token required",
			})
			ctx.Abort()
			return
		}

		token := texts[1]

		claims, err := m.authService.AuthenticateToken(token)
		if err != nil {
			m.logger.Errorf("authenticating token failed: %v", err)

			ctx.JSON(401, gin.H{
				"error": "Invalid token",
			})
			ctx.Abort()
			return
		}

		ctx.Set(constants.JWTClaims, claims.TokenPayload)
		ctx.Next()
	}
}
