package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/dtos"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
	"github.com/shunwuse/go-hris/services"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	logger      lib.Logger
	userService services.UserService
	authService services.AuthService
}

func NewUserController() UserController {
	logger := lib.GetLogger()

	// Initialize services
	userService := services.NewUserService()
	authService := services.NewAuthService()

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
	token := ctx.MustGet(constants.JWTClaims).(services.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionReadUser); !hasPermission {
		c.logger.Errorf("Error user not authorized to get users")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to get users",
		})
		return
	}

	users, err := c.userService.GetUsers()
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
	var userCreate dtos.UserCreate

	token := ctx.MustGet(constants.JWTClaims).(services.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionCreateUser); !hasPermission {
		c.logger.Errorf("Error user not authorized to create user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to create user",
		})
		return
	}

	if err := ctx.ShouldBindJSON(&userCreate); err != nil {
		c.logger.Errorf("Error binding user: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// cannot create user with role admin
	if userCreate.Role == constants.Admin {
		c.logger.Errorf("Error user not authorized to create admin user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to create admin user",
		})
		return
	}

	hashedPassword, err := hashPassword(userCreate.Password)
	if err != nil {
		c.logger.Errorf("Error hashing password: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error hashing password",
		})
	}

	var user = &models.User{
		Username: userCreate.Username,
		Name:     userCreate.Name,

		Password: models.Password{
			Hash: hashedPassword,
		},
	}

	if err := c.userService.CreateUser(user, userCreate.Role); err != nil {
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
	var userUpdate dtos.UserUpdate

	token := ctx.MustGet(constants.JWTClaims).(services.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionUpdateUser); !hasPermission {
		c.logger.Errorf("Error user not authorized to update user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to update user",
		})
		return
	}

	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		c.logger.Errorf("Error binding user: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	var user = &models.User{
		ID:   userUpdate.ID,
		Name: userUpdate.Name,
	}

	if err := c.userService.UpdateUser(user); err != nil {
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

	user, err := c.userService.GetUserByUsername(userLogin.Username)
	if err != nil {
		c.logger.Errorf("Error getting user(%s): %v", userLogin.Username, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "user not found",
		})
		return
	}

	// check password
	passwordMatch := checkPasswordHash(userLogin.Password, user.Password.Hash)
	if !passwordMatch {
		c.logger.Errorf("Error password not match")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Password not match",
		})
		return
	}

	// generate token
	token, err := c.authService.GenerateToken(user)
	if err != nil {
		c.logger.Errorf("generating token failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error generating token",
		})
		return
	}

	roles := make([]string, 0)
	for _, role := range user.Roles {
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

func hashPassword(password string) (string, error) {
	// Use bcrypt to hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // if err is nil, password match
}
