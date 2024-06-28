package services

import (
	"errors"
	"slices"

	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
	"github.com/shunwuse/go-hris/repositories"
)

type UserService struct {
	logger                   lib.Logger
	userRepository           repositories.UserRepository
	roleRepository           repositories.RoleRepository
	userRoleRepository       repositories.UserRoleRepository
	rolePermissionRepository repositories.RolePermissionRepository
}

func NewUserService() UserService {
	logger := lib.GetLogger()

	// Initialize repositories
	userRepository := repositories.NewUserRepository()
	roleRepository := repositories.NewRoleRepository()
	userRoleRepository := repositories.NewUserRoleRepository()
	rolePermissionRepository := repositories.NewRolePermissionRepository()

	return UserService{
		logger:                   logger,
		userRepository:           userRepository,
		roleRepository:           roleRepository,
		userRoleRepository:       userRoleRepository,
		rolePermissionRepository: rolePermissionRepository,
	}
}

func (s UserService) GetUsers() ([]models.User, error) {
	var users []models.User

	result := s.userRepository.Find(&users)
	if result.Error != nil {
		s.logger.Errorf("Error getting users: %v", result.Error)
		return nil, result.Error
	}

	return users, nil
}

func (s UserService) CreateUser(user *models.User, role constants.Role) error {
	result := s.userRepository.Create(user)
	if result.Error != nil {
		s.logger.Errorf("Error creating user: %v", result.Error)
		return result.Error
	}

	roleModel := s.roleRepository.GetRoleByName(role.String())
	if roleModel == nil {
		// s.logger.Infof("Role not found, creating role: %v", role)

		// roleModel = &models.Role{
		// 	Name: constants.Staff.String(),
		// }

		// if err := s.roleRepository.AddRole(roleModel); err != nil {
		// 	s.logger.Errorf("add role error: %v", err)
		// 	return err
		// }

		s.logger.Errorf("Role not found: %v", role)
		return errors.New("role not found")
	}

	// Add user role
	userRole := &models.UserRole{
		UserID: user.ID,
		RoleID: roleModel.ID,
	}

	// Create user role
	result = s.userRoleRepository.Create(userRole)
	if result.Error != nil {
		s.logger.Errorf("creating user role error: %v", result.Error)
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

	// Get permissions
	permissions := make(constants.Permissions, 0)
	roles := user.Roles

	// Get permissions by role
	for _, role := range roles {
		rolePermissions := s.rolePermissionRepository.GetPermissionsByRole(constants.Role(role.Name))

		// Add permissions to user
		for _, permission := range rolePermissions {
			if !slices.Contains(permissions, permission) {
				permissions = append(permissions, permission)
			}
		}
	}

	// Set permissions to user
	user.Permissions = permissions

	return user, nil
}

func (s UserService) UpdateUser(user *models.User) error {
	result := s.userRepository.Updates(user)
	if result.Error != nil {
		s.logger.Errorf("Error updating user: %v", result.Error)
		return result.Error
	}

	return nil
}
