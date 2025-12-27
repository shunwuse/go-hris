package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type UserController struct {
	logger      *infra.Logger
	userService service.UserService
}

func NewUserController(
	logger *infra.Logger,
	userService service.UserService,
) *UserController {
	return &UserController{
		logger:      logger,
		userService: userService,
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

	var user = &domains.UserCreate{
		Username: userCreate.Username,
		Name:     userCreate.Name,
		Password: userCreate.Password,
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
