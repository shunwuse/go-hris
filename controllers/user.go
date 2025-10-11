package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/domains"
	"github.com/shunwuse/go-hris/dtos"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/lib/utils"
	"github.com/shunwuse/go-hris/ports/service"
)

type UserController struct {
	logger      lib.Logger
	userService service.UserService
	authService service.AuthService
}

func NewUserController(
	logger lib.Logger,
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
func (c UserController) GetUsers(ctx *gin.Context) {
	token := ctx.MustGet(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionReadUser); !hasPermission {
		c.logger.Errorf("Error user not authorized to get users")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to get users",
		})
		return
	}

	users, err := c.userService.GetUsers(ctx)
	if err != nil {
		c.logger.Errorf("Error getting users: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
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

	ctx.JSON(http.StatusOK, gin.H{
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
func (c UserController) CreateUser(ctx *gin.Context) {
	token := ctx.MustGet(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionCreateUser); !hasPermission {
		c.logger.Errorf("Error user not authorized to create user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to create user",
		})
		return
	}

	var userCreate dtos.UserCreate
	if err := ctx.ShouldBindJSON(&userCreate); err != nil {
		c.logger.Errorf("Error binding user: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// lower case username
	userCreate.Username = strings.ToLower(userCreate.Username)

	// cannot create user with role admin
	if userCreate.Role == constants.Admin {
		c.logger.Errorf("Error user not authorized to create admin user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to create admin user",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(userCreate.Password)
	if err != nil {
		c.logger.Errorf("Error hashing password: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
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

	if err := c.userService.CreateUser(ctx, user, userCreate.Role); err != nil {
		c.logger.Errorf("Error creating user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
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
func (c UserController) UpdateUser(ctx *gin.Context) {
	token := ctx.MustGet(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionUpdateUser); !hasPermission {
		c.logger.Errorf("Error user not authorized to update user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to update user",
		})
		return
	}

	var userUpdate dtos.UserUpdate
	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		c.logger.Errorf("Error binding user: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	var user = &domains.UserUpdate{
		ID:   userUpdate.ID,
		Name: userUpdate.Name,
	}

	if err := c.userService.UpdateUser(ctx, user); err != nil {
		c.logger.Errorf("Error updating user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error updating user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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
func (c UserController) Login(ctx *gin.Context) {
	var userLogin dtos.UserLogin

	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		c.logger.Errorf("Error binding user: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// lower case username
	userLogin.Username = strings.ToLower(userLogin.Username)

	user, err := c.userService.GetUserByUsername(ctx, userLogin.Username)
	if err != nil {
		c.logger.Errorf("Error getting user(%s): %v", userLogin.Username, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "user not found",
		})
		return
	}

	// check password
	passwordMatch := utils.CheckPasswordHash(userLogin.Password, user.Edges.Password.Hash)
	if !passwordMatch {
		c.logger.Errorf("Error password not match")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Password not match",
		})
		return
	}

	// generate token
	token, err := c.authService.GenerateToken(ctx, user)
	if err != nil {
		c.logger.Errorf("generating token failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
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

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
