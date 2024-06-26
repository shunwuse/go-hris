package services

import (
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
	"github.com/shunwuse/go-hris/repositories"
)

type UserService struct {
	logger         lib.Logger
	userRepository repositories.UserRepository
}

func NewUserService() UserService {
	logger := lib.GetLogger()

	// Initialize repositories
	userRepository := repositories.NewUserRepository()

	return UserService{
		logger:         logger,
		userRepository: userRepository,
	}
}

func (s UserService) CreateUser(user *models.User) error {
	result := s.userRepository.Create(user)
	if result.Error != nil {
		s.logger.Errorf("Error creating user: %v", result.Error)
		return result.Error
	}

	return nil
}

func (s UserService) GetUserByUsername(username string) (*models.User, error) {
	var user *models.User

	result := s.userRepository.Preload("Password").Preload("Roles").First(&user, "username = ?", username)
	if result.Error != nil {
		s.logger.Errorf("Error getting user by username: %v", result.Error)
		return nil, result.Error
	}

	return user, nil
}
