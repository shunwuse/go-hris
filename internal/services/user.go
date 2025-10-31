package services

import (
	"context"
	"slices"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
	"go.uber.org/zap"
)

type userService struct {
	logger                   *infra.Logger
	userRepository           *repositories.UserRepository
	passwordRepository       *repositories.PasswordRepository
	roleRepository           *repositories.RoleRepository
	userRoleRepository       *repositories.UserRoleRepository
	rolePermissionRepository *repositories.RolePermissionRepository
}

func NewUserService(
	logger *infra.Logger,
	userRepository *repositories.UserRepository,
	passwordRepository *repositories.PasswordRepository,
	roleRepository *repositories.RoleRepository,
	userRoleRepository *repositories.UserRoleRepository,
	rolePermissionRepository *repositories.RolePermissionRepository,
) service.UserService {
	return &userService{
		logger:                   logger,
		userRepository:           userRepository,
		passwordRepository:       passwordRepository,
		roleRepository:           roleRepository,
		userRoleRepository:       userRoleRepository,
		rolePermissionRepository: rolePermissionRepository,
	}
}

func (s *userService) GetUsers(ctx context.Context) ([]*entgen.User, error) {
	return s.userRepository.FindAll(ctx)
}

func (s *userService) CreateUser(ctx context.Context, user *domains.UserCreate, role constants.Role) error {
	u, err := s.userRepository.Create(ctx, user.Username, user.Name)
	if err != nil {
		return err
	}

	_, err = s.passwordRepository.Create(ctx, user.Password.Hash, u)
	if err != nil {
		return err
	}

	roleModel := s.roleRepository.FindByName(ctx, role.String())
	if roleModel == nil {
		s.logger.WithContext(ctx).Error("role not found", zap.String("role", role.String()))
		return errors.ErrNotFound
	}

	_, err = s.userRoleRepository.Create(ctx, u.ID, roleModel.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error) {
	user, err := s.userRepository.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Get permissions based on user's roles.
	permissions := make(constants.Permissions, 0)

	// Collect permissions from all roles.
	for _, role := range user.Edges.Roles {
		rolePermissions := s.rolePermissionRepository.GetPermissionsByRole(ctx, constants.Role(role.Name))

		// Add unique permissions to user.
		for _, permission := range rolePermissions {
			if !slices.Contains(permissions, permission) {
				permissions = append(permissions, permission)
			}
		}
	}

	// Construct user with permissions.
	u := domains.UserWithPermissions{
		User:        user,
		Permissions: permissions,
	}

	return &u, nil
}

func (s *userService) UpdateUser(ctx context.Context, update *domains.UserUpdate) error {
	return s.userRepository.Update(ctx, update.ID, update.Name)
}
