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

// GetUsers controller
func (c UserController) GetUsers(ctx *gin.Context) {
	token := ctx.MustGet(constants.JWTClaims).(services.TokenPayload)
	roles := token.Roles

	// check all roles
	hasAdminRole := false
	for _, role := range roles {
		hasAdminRole = role.IsAdmin()
	}

	if !hasAdminRole {
		c.logger.Errorf("Error user not authorized")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized",
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

// CreateUser controller
func (c UserController) CreateUser(ctx *gin.Context) {
	var userCreate dtos.UserCreate

	token := ctx.MustGet(constants.JWTClaims).(services.TokenPayload)
	roles := token.Roles

	// check all roles
	hasAdminRole := false
	for _, role := range roles {
		hasAdminRole = role.IsAdmin()
	}

	if !hasAdminRole {
		c.logger.Errorf("Error user not authorized")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized",
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

	if err := c.userService.CreateUser(user); err != nil {
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

// UpdateUser controller
func (c UserController) UpdateUser(ctx *gin.Context) {
	var userUpdate dtos.UserUpdate

	token := ctx.MustGet(constants.JWTClaims).(services.TokenPayload)
	roles := token.Roles

	// check all roles
	hasAdminRole := false
	for _, role := range roles {
		hasAdminRole = role.IsAdmin()
	}

	if !hasAdminRole {
		c.logger.Errorf("Error user not authorized")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized",
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

}

// Login controller
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
