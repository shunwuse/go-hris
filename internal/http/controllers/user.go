package controllers

import (
	"net/http"
	"strconv"

	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/request"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type UserController struct {
	logger      *logger.Logger
	userService service.UserService
}

func NewUserController(
	log *logger.Logger,
	userService service.UserService,
) *UserController {
	return &UserController{
		logger:      log,
		userService: userService,
	}
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Get pagination and filter parameters.
	var query dtos.UserGetList
	if err := request.DecodeQuery(r, &query); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode query params", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}
	query.Normalize()

	if err := query.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate query params", zap.Error(err))
		response.Error(w, err)
		return
	}

	filter := domains.UserFilter{
		Name: query.Name,
		Role: query.Role,
	}

	result, err := c.userService.GetUsersWithOffset(r.Context(), query.OffsetQuery, filter)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to get users", zap.Error(err))
		response.Error(w, err)
		return
	}

	usersResponse := make([]dtos.UserResponse, len(result.Items))
	for idx, user := range result.Items {
		usersResponse[idx] = dtos.UserResponse{
			ID:              user.ID,
			Username:        user.Username,
			Name:            user.Name,
			CreatedTime:     strconv.FormatInt(user.CreatedAt.UnixMilli(), 10),
			LastUpdatedTime: strconv.FormatInt(user.UpdatedAt.UnixMilli(), 10),
		}
	}

	response.OffsetList(w, usersResponse, response.OffsetPaginationMeta{
		Total:       result.TotalCount,
		PerPage:     query.PerPage,
		CurrentPage: query.Page,
		LastPage:    result.TotalPage,
	})
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	var pathParams dtos.UserPathParams
	if err := request.DecodePath(r, &pathParams); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode path params", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if err := pathParams.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate path params", zap.Error(err))
		response.Error(w, err)
		return
	}

	user, err := c.userService.GetUserByID(r.Context(), pathParams.ID)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to get user", zap.Error(err), zap.Uint("user_id", pathParams.ID))
		response.Error(w, err)
		return
	}

	resp := dtos.UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		Name:            user.Name,
		CreatedTime:     strconv.FormatInt(user.CreatedAt.UnixMilli(), 10),
		LastUpdatedTime: strconv.FormatInt(user.UpdatedAt.UnixMilli(), 10),
	}

	response.OK(w, resp)
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userCreate dtos.UserCreate
	if err := request.DecodeJSON(r, &userCreate); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode user request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if err := userCreate.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate user request", zap.Error(err))
		response.Error(w, err)
		return
	}

	var user = &domains.UserCreate{
		Username: userCreate.Username,
		Name:     userCreate.Name,
		Password: userCreate.Password,
	}

	err := c.userService.CreateUser(r.Context(), user, userCreate.Role)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to create user", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Created(w, "user created successfully")
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var pathParams dtos.UserPathParams
	if err := request.DecodePath(r, &pathParams); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode path params", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if err := pathParams.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate path params", zap.Error(err))
		response.Error(w, err)
		return
	}

	var userUpdate dtos.UserUpdate
	if err := request.DecodeJSON(r, &userUpdate); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode user request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if err := userUpdate.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate user request", zap.Error(err))
		response.Error(w, err)
		return
	}

	var user = &domains.UserUpdate{
		ID:   pathParams.ID,
		Name: userUpdate.Name,
	}

	err := c.userService.UpdateUser(r.Context(), user)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to update user", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.OK(w, "user updated successfully")
}
