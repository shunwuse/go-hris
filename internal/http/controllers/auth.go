package controllers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
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

func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dtos.RefreshRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode refresh request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if req.RefreshToken == "" {
		c.logger.WithContext(r.Context()).Error("refresh token is required")
		response.Error(w, errors.ErrInvalidInput)
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

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode logout request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if req.RefreshToken == "" {
		c.logger.WithContext(r.Context()).Error("refresh token is required")
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	// Revoke refresh token.
	if err := c.authService.RevokeRefreshToken(r.Context(), req.RefreshToken); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to revoke refresh token", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.OK(w, "logged out successfully")
}
