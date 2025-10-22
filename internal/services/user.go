package services

import (
	"context"
	"errors"
	"slices"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/user"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
	"go.uber.org/zap"
)

type userService struct {
	logger                   infra.Logger
	userRepository           repositories.UserRepository
	roleRepository           repositories.RoleRepository
	userRoleRepository       repositories.UserRoleRepository
	rolePermissionRepository repositories.RolePermissionRepository
}

func NewUserService(
	logger infra.Logger,
	userRepository repositories.UserRepository,
	roleRepository repositories.RoleRepository,
	userRoleRepository repositories.UserRoleRepository,
	rolePermissionRepository repositories.RolePermissionRepository,
) service.UserService {
	return userService{
		logger:                   logger,
		userRepository:           userRepository,
		roleRepository:           roleRepository,
		userRoleRepository:       userRoleRepository,
		rolePermissionRepository: rolePermissionRepository,
	}
}

func (s userService) GetUsers(ctx context.Context) ([]*entgen.User, error) {
	users, err := s.userRepository.Client.User.
		Query().
		All(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("Error getting users", zap.Error(err))
		return nil, err
	}

	return users, nil
}

func (s userService) CreateUser(ctx context.Context, user *domains.UserCreate, role constants.Role) error {
	u, err := s.userRepository.Client.User.
		Create().
		SetUsername(user.Username).
		SetName(user.Name).
		Save(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("Error creating user", zap.Error(err))
		return err
	}

	_, err = s.userRepository.Client.Password.
		Create().
		SetHash(user.Password.Hash).
		SetOwner(u).
		Save(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("Error creating password", zap.Error(err))
		return err
	}

	roleModel := s.roleRepository.GetRoleByName(ctx, role.String())
	if roleModel == nil {
		// s.logger.WithContext(ctx).Info("Role not found, creating role", zap.String("role", role.String()))

		// roleCreate := &domains.RoleCreate{
		// 	Name: constants.Staff.String(),
		// }

		// if err := s.roleRepository.AddRole(ctx, roleCreate); err != nil {
		// 	s.logger.WithContext(ctx).Error("add role error", zap.Error(err))
		// 	return err
		// }

		s.logger.WithContext(ctx).Error("Role not found", zap.String("role", role.String()))
		return errors.New("role not found")
	}

	// Create user role
	_, err = s.userRepository.Client.UserRole.
		Create().
		SetUserID(u.ID).
		SetRoleID(roleModel.ID).
		Save(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("creating user role error", zap.Error(err))
		return err
	}

	return nil
}

func (s userService) GetUserByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error) {
	user, err := s.userRepository.Client.User.
		Query().
		WithPassword().
		WithRoles().
		Where(user.UsernameEQ(username)).
		Only(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("Error getting user by username", zap.Error(err))
		return nil, err
	}

	// Get permissions
	permissions := make(constants.Permissions, 0)
	// roles := user.Roles

	// Get permissions by role
	for _, role := range user.Edges.Roles {
		rolePermissions := s.rolePermissionRepository.GetPermissionsByRole(ctx, constants.Role(role.Name))

		// Add permissions to user
		for _, permission := range rolePermissions {
			if !slices.Contains(permissions, permission) {
				permissions = append(permissions, permission)
			}
		}
	}

	// Set permissions to user
	u := domains.UserWithPermissions{
		User:        user,
		Permissions: permissions,
	}

	return &u, nil
}

func (s userService) UpdateUser(ctx context.Context, update *domains.UserUpdate) error {
	err := s.userRepository.Client.User.
		Update().
		Where(user.IDEQ(update.ID)).
		SetName(update.Name).
		Exec(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("Error updating user", zap.Error(err))
		return err
	}

	return nil
}
