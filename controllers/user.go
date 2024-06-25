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
}

func NewUserController() UserController {
	logger := lib.GetLogger()

	// Initialize services
	userService := services.NewUserService()

	return UserController{
		logger:      logger,
		userService: userService,
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
		Username:     userCreate.Username,
		PasswordHash: hashedPassword,
		Name:         userCreate.Name,
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

func hashPassword(password string) (string, error) {
	// Use bcrypt to hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
