package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// CreateUser controller
func (c UserController) CreateUser(ctx *gin.Context) {
	var userCreate dtos.UserCreate

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

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
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
