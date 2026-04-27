package controllers

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/request"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type AuthController struct {
	logger      *logger.Logger
	authService service.AuthService
}

func NewAuthController(
	log *logger.Logger,
	authService service.AuthService,
) *AuthController {
	return &AuthController{
		logger:      log,
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

	claims, _ := contextx.GetClaims(r.Context())
	if err := c.authService.Logout(r.Context(), req.RefreshToken, claims); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to logout", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.OK(w, "logged out successfully")
}

func (c *AuthController) LogoutAll(w http.ResponseWriter, r *http.Request) {
	claims, _ := contextx.GetClaims(r.Context())
	if err := c.authService.LogoutAll(r.Context(), claims); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to logout all", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.OK(w, "logged out from all devices successfully")
}
