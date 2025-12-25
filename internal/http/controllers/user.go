package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/utils"
	"go.uber.org/zap"
)

type UserController struct {
	logger      *infra.Logger
	userService service.UserService
	authService service.AuthService
}

func NewUserController(
	logger *infra.Logger,
	userService service.UserService,
	authService service.AuthService,
) *UserController {
	return &UserController{
		logger:      logger,
		userService: userService,
		authService: authService,
	}
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// Check if user has permission to read users.
	if hasPermission := permissions.Contains(constants.PermissionReadUser); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to get users")
		response.Error(w, errors.ErrInsufficientPermissions)
		return
	}

	users, err := c.userService.GetUsers(r.Context())
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to get users", zap.Error(err))
		response.Error(w, err)
		return
	}

	usersResponse := make([]dtos.GetUserResponse, 0)
	for _, user := range users {
		usersResponse = append(usersResponse, dtos.GetUserResponse{
			ID:              user.ID,
			Username:        user.Username,
			Name:            user.Name,
			CreatedTime:     strconv.FormatInt(user.CreatedAt.UnixMilli(), 10),
			LastUpdatedTime: strconv.FormatInt(user.UpdatedAt.UnixMilli(), 10),
		})
	}

	response.List(w, usersResponse)
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// Check if user has permission to create users.
	if hasPermission := permissions.Contains(constants.PermissionCreateUser); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to create user")
		response.Error(w, errors.ErrInsufficientPermissions)
		return
	}

	var userCreate dtos.UserCreate
	if err := render.DecodeJSON(r.Body, &userCreate); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode user request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	// Convert username to lowercase.
	userCreate.Username = strings.ToLower(userCreate.Username)

	// Cannot create user with admin role.
	if userCreate.Role == constants.Admin {
		c.logger.WithContext(r.Context()).Error("cannot create user with admin role")
		response.Error(w, errors.ErrOperationNotAllowed)
		return
	}

	hashedPassword, err := utils.HashPassword(userCreate.Password)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to hash password", zap.Error(err))
		response.Error(w, errors.ErrInternalError)
		return
	}

	var user = &domains.UserCreate{
		Username: userCreate.Username,
		Name:     userCreate.Name,

		Password: domains.PasswordCreate{
			Hash: hashedPassword,
		},
	}

	if err := c.userService.CreateUser(r.Context(), user, userCreate.Role); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to create user", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Created(w, "user created successfully")
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// Check if user has permission to update users.
	if hasPermission := permissions.Contains(constants.PermissionUpdateUser); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to update user")
		response.Error(w, errors.ErrInsufficientPermissions)
		return
	}

	var userUpdate dtos.UserUpdate
	if err := render.DecodeJSON(r.Body, &userUpdate); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode user request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	var user = &domains.UserUpdate{
		ID:   userUpdate.ID,
		Name: userUpdate.Name,
	}

	if err := c.userService.UpdateUser(r.Context(), user); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to update user", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.OK(w, "user updated successfully")
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var userLogin dtos.UserLogin

	if err := render.DecodeJSON(r.Body, &userLogin); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode login request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	// Convert username to lowercase.
	userLogin.Username = strings.ToLower(userLogin.Username)

	user, err := c.userService.GetUserByUsername(r.Context(), userLogin.Username)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to get user", zap.String("username", userLogin.Username), zap.Error(err))
		response.Error(w, err)
		return
	}

	// Verify password.
	passwordMatch := utils.CheckPasswordHash(userLogin.Password, user.Edges.Password.Hash)
	if !passwordMatch {
		c.logger.WithContext(r.Context()).Error("password verification failed")
		response.Error(w, errors.ErrInvalidCredentials)
		return
	}

	// Generate access token.
	accessToken, err := c.authService.GenerateAccessToken(r.Context(), user)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to generate access token", zap.Error(err))
		response.Error(w, err)
		return
	}

	// Generate refresh token.
	refreshToken, err := c.authService.GenerateRefreshToken(r.Context(), user)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to generate refresh token", zap.Error(err))
		response.Error(w, err)
		return
	}

	roles := make([]constants.Role, 0)
	for _, role := range user.Edges.Roles {
		roles = append(roles, role.Name)
	}

	resp := dtos.LoginResponse{
		Username:     user.Username,
		Roles:        roles,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.OK(w, resp)
}
