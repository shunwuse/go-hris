package services

import (
	"context"
	"slices"
	"strings"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
	"github.com/shunwuse/go-hris/internal/utils"
	"go.uber.org/zap"
)

type userService struct {
	logger         *infra.Logger
	userRepository *repositories.UserRepository
	roleRepository *repositories.RoleRepository
}

func NewUserService(
	logger *infra.Logger,
	userRepository *repositories.UserRepository,
	roleRepository *repositories.RoleRepository,
) service.UserService {
	return &userService{
		logger:         logger,
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
}

func (s *userService) GetUsers(ctx context.Context) ([]*entgen.User, error) {
	return s.userRepository.FindAll(ctx)
}

func (s *userService) GetUsersWithOffset(ctx context.Context, query domains.OffsetQuery) (*domains.OffsetResult[*entgen.User], error) {
	return s.userRepository.FindAllWithOffset(ctx, query)
}

func (s *userService) CreateUser(ctx context.Context, user *domains.UserCreate, role constants.Role) error {
	// Convert username to lowercase.
	user.Username = strings.ToLower(user.Username)

	// Cannot create user with admin role.
	if role == constants.Admin {
		s.logger.WithContext(ctx).Error("cannot create user with admin role")
		return errors.ErrOperationNotAllowed
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to hash password", zap.Error(err))
		return errors.ErrInternalError
	}

	return s.userRepository.WithTx(ctx, func(txCtx context.Context) error {
		u, err := s.userRepository.Create(txCtx, user.Username, user.Name)
		if err != nil {
			return err
		}

		_, err = s.userRepository.CreatePassword(txCtx, hashedPassword, u)
		if err != nil {
			return err
		}

		roleModel := s.roleRepository.FindByName(txCtx, role)
		if roleModel == nil {
			s.logger.WithContext(txCtx).Error("role not found", zap.String("role", role.String()))
			return errors.ErrNotFound
		}

		_, err = s.userRepository.AssignRole(txCtx, u.ID, roleModel.ID)
		if err != nil {
			return err
		}

		return nil
	})
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
		rolePermissions := s.roleRepository.GetPermissionsByRole(ctx, constants.Role(role.Name))

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

func (s *userService) GetUserByID(ctx context.Context, id uint) (*domains.UserWithPermissions, error) {
	return s.userRepository.FindByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, update *domains.UserUpdate) error {
	return s.userRepository.Update(ctx, update.ID, update.Name)
}
