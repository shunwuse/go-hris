package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/utils"
	"go.uber.org/zap"
)

type UserController struct {
	logger      infra.Logger
	userService service.UserService
	authService service.AuthService
}

func NewUserController(
	logger infra.Logger,
	userService service.UserService,
	authService service.AuthService,
) UserController {
	return UserController{
		logger:      logger,
		userService: userService,
		authService: authService,
	}
}

// GetUsers godoc
//
// @Summary Get users
// @Description Get all users
// @Tags users
// @security BasicAuth
// @Accept json
// @Produce json
// @Success 200 {array} dtos.GetUserResponse
// @Router /users [get]
func (c UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionReadUser); !hasPermission {
		c.logger.Error("Error user not authorized to get users")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "User not authorized to get users",
		})
		return
	}

	users, err := c.userService.GetUsers(r.Context())
	if err != nil {
		c.logger.Error("Error getting users", zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Error getting users",
		})
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

	render.JSON(w, r, map[string]any{
		"data": usersResponse,
	})
}

// CreateUser godoc
//
// @Summary Create user
// @Description Create a new user
// @Tags users
// @security BasicAuth
// @Accept json
// @Produce json
// @Param user body dtos.UserCreate true "User object that needs to be created"
// @Success 201 {string} string "create successfully"
// @Router /users [post]
func (c UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionCreateUser); !hasPermission {
		c.logger.Error("Error user not authorized to create user")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "User not authorized to create user",
		})
		return
	}

	var userCreate dtos.UserCreate
	if err := render.DecodeJSON(r.Body, &userCreate); err != nil {
		c.logger.Error("Error binding user", zap.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Invalid request",
		})
		return
	}

	// lower case username
	userCreate.Username = strings.ToLower(userCreate.Username)

	// cannot create user with role admin
	if userCreate.Role == constants.Admin {
		c.logger.Error("Error user not authorized to create admin user")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "User not authorized to create admin user",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(userCreate.Password)
	if err != nil {
		c.logger.Error("Error hashing password", zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Error hashing password",
		})
	}

	var user = &domains.UserCreate{
		Username: userCreate.Username,
		Name:     userCreate.Name,

		Password: domains.PasswordCreate{
			Hash: hashedPassword,
		},
	}

	if err := c.userService.CreateUser(r.Context(), user, userCreate.Role); err != nil {
		c.logger.Error("Error creating user", zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Error creating user",
		})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{
		"message": "create successfully",
	})
}

// UpdateUser godoc
//
// @Summary Update user
// @Description Update user
// @Tags users
// @security BasicAuth
// @Accept json
// @Produce json
// @Param user body dtos.UserUpdate true "User object that needs to be updated"
// @Success 200 {string} string "update successfully"
// @Router /users [put]
func (c UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionUpdateUser); !hasPermission {
		c.logger.Error("Error user not authorized to update user")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "User not authorized to update user",
		})
		return
	}

	var userUpdate dtos.UserUpdate
	if err := render.DecodeJSON(r.Body, &userUpdate); err != nil {
		c.logger.Error("Error binding user", zap.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Invalid request",
		})
		return
	}

	var user = &domains.UserUpdate{
		ID:   userUpdate.ID,
		Name: userUpdate.Name,
	}

	if err := c.userService.UpdateUser(r.Context(), user); err != nil {
		c.logger.Error("Error updating user", zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Error updating user",
		})
		return
	}

	render.JSON(w, r, map[string]string{
		"message": "update successfully",
	})
}

// Login godoc
//
// @Summary Login
// @Description Login
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.UserLogin true "User object that needs to login"
// @Success 200 {object} dtos.LoginResponse
// @Router /login [post]
func (c UserController) Login(w http.ResponseWriter, r *http.Request) {
	var userLogin dtos.UserLogin

	if err := render.DecodeJSON(r.Body, &userLogin); err != nil {
		c.logger.Error("Error binding user", zap.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Invalid request",
		})
		return
	}

	// lower case username
	userLogin.Username = strings.ToLower(userLogin.Username)

	user, err := c.userService.GetUserByUsername(r.Context(), userLogin.Username)
	if err != nil {
		c.logger.Error("error getting user", zap.String("username", userLogin.Username), zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "user not found",
		})
		return
	}

	// check password
	passwordMatch := utils.CheckPasswordHash(userLogin.Password, user.Edges.Password.Hash)
	if !passwordMatch {
		c.logger.Error("Error password not match")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "Password not match",
		})
		return
	}

	// generate token
	token, err := c.authService.GenerateToken(r.Context(), user)
	if err != nil {
		c.logger.Error("generating token failed", zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Error generating token",
		})
		return
	}

	roles := make([]string, 0)
	for _, role := range user.Edges.Roles {
		roles = append(roles, role.Name)
	}

	response := dtos.LoginResponse{
		Username: user.Username,
		Roles:    roles,
		Token:    token,
	}

	render.JSON(w, r, map[string]any{
		"data": response,
	})
}
