package controllers

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/request"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type AuthController struct {
	logger      *infra.Logger
	authService service.AuthService
}

func NewAuthController(
	logger *infra.Logger,
	authService service.AuthService,
) *AuthController {
	return &AuthController{
		logger:      logger,
		authService: authService,
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var userLogin dtos.UserLogin

	if err := request.DecodeJSON(r, &userLogin); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode login request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if err := userLogin.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate login request", zap.Error(err))
		response.Error(w, err)
		return
	}

	result, err := c.authService.Login(r.Context(), userLogin.Username, userLogin.Password)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to login", zap.String("username", userLogin.Username), zap.Error(err))
		response.Error(w, err)
		return
	}

	resp := dtos.UserLoginResponse{
		Username:     result.Username,
		Roles:        result.Roles,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	response.OK(w, resp)
}

func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dtos.RefreshRequest

	if err := request.DecodeJSON(r, &req); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode refresh request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if err := req.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate refresh request", zap.Error(err))
		response.Error(w, err)
		return
	}

	tokenPair, err := c.authService.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to refresh tokens", zap.Error(err))
		response.Error(w, err)
		return
	}

	resp := dtos.RefreshResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	response.OK(w, resp)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	var req dtos.LogoutRequest

	_ = request.DecodeJSON(r, &req) // Ignore error, as refresh token is optional

	if req.RefreshToken != "" {
		// Revoke refresh token.
		if err := c.authService.RevokeRefreshToken(r.Context(), req.RefreshToken); err != nil {
			c.logger.WithContext(r.Context()).Error("failed to revoke refresh token", zap.Error(err))
		}
	}

	// Blacklist current access token.
	claims, ok := request.GetClaims(r)
	if ok && claims.JTI != "" {
		expiration := claims.ExpiresIn()
		if expiration > 0 {
			if err := c.authService.BlacklistToken(r.Context(), claims.JTI, expiration); err != nil {
				c.logger.WithContext(r.Context()).Error("failed to blacklist token", zap.Error(err))
			}
		}
	}

	response.OK(w, "logged out successfully")
}

func (c *AuthController) LogoutAll(w http.ResponseWriter, r *http.Request) {
	claims, ok := request.GetClaims(r)
	if !ok {
		c.logger.WithContext(r.Context()).Error("failed to get claims from context")
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	// Revoke all refresh tokens for this user.
	if err := c.authService.RevokeAllUserTokens(r.Context(), claims.Identity.UserID); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to revoke all tokens", zap.Uint("user_id", claims.Identity.UserID), zap.Error(err))
		response.Error(w, err)
		return
	}

	// Blacklist current access token.
	if claims.JTI != "" {
		expiration := claims.ExpiresIn()
		if expiration > 0 {
			if err := c.authService.BlacklistToken(r.Context(), claims.JTI, expiration); err != nil {
				c.logger.WithContext(r.Context()).Error("failed to blacklist token", zap.Error(err))
			}
		}
	}

	response.OK(w, "logged out from all devices successfully")
}
