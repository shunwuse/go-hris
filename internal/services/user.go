package services

import (
	"context"
	"strings"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/pkg/cryptox"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/infra"
	"github.com/shunwuse/go-hris/internal/ports/query"
	"github.com/shunwuse/go-hris/internal/ports/repository"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type userService struct {
	logger     *logger.Logger
	transactor infra.Transactor

	userRepository repository.UserRepository
	roleRepository repository.RoleRepository

	reader query.UserIdentityReader
}

func NewUserService(
	log *logger.Logger,
	transactor infra.Transactor,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	reader query.UserIdentityReader,
) service.UserService {
	return &userService{
		logger:         log,
		transactor:     transactor,
		userRepository: userRepository,
		roleRepository: roleRepository,
		reader:         reader,
	}
}

func (s *userService) GetUsers(ctx context.Context) ([]*entgen.User, error) {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return nil, errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionReadUser) {
		s.logger.WithContext(ctx).Error("not authorized to get users")
		return nil, errors.ErrForbidden
	}

	return s.userRepository.FindAll(ctx)
}

func (s *userService) GetUsersWithOffset(ctx context.Context, query domains.OffsetQuery, filter domains.UserFilter) (*domains.OffsetResult[*entgen.User], error) {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return nil, errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionReadUser) {
		s.logger.WithContext(ctx).Error("not authorized to get users")
		return nil, errors.ErrForbidden
	}

	return s.userRepository.FindAllWithOffset(ctx, query, filter)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error) {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return nil, errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionReadUser) {
		s.logger.WithContext(ctx).Error("not authorized to get user")
		return nil, errors.ErrForbidden
	}

	return s.reader.GetUserWithPermissionsByUsername(ctx, username)
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*domains.UserWithPermissions, error) {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return nil, errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionReadUser) {
		s.logger.WithContext(ctx).Error("not authorized to get user")
		return nil, errors.ErrForbidden
	}

	return s.reader.GetUserWithPermissionsByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, user *domains.UserCreate, role constants.Role) error {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionCreateUser) {
		s.logger.WithContext(ctx).Error("not authorized to create user")
		return errors.ErrForbidden
	}

	// Convert username to lowercase.
	user.Username = strings.ToLower(user.Username)

	// Cannot create user with admin role.
	if role == constants.Admin {
		s.logger.WithContext(ctx).Error("cannot create user with admin role")
		return errors.ErrOperationNotAllowed
	}

	hashedPassword, err := cryptox.HashPassword(user.Password)
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to hash password", zap.Error(err))
		return errors.ErrInternalError
	}

	return s.transactor.WithTx(ctx, func(txCtx context.Context) error {
		u, err := s.userRepository.Create(txCtx, user.Username, user.Name)
		if err != nil {
			return err
		}

		_, err = s.userRepository.CreatePassword(txCtx, hashedPassword, u)
		if err != nil {
			return err
		}

		roleModel := s.roleRepository.FindRoleByName(txCtx, role)
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

func (s *userService) UpdateUser(ctx context.Context, update *domains.UserUpdate) error {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionUpdateUser) {
		s.logger.WithContext(ctx).Error("not authorized to update user")
		return errors.ErrForbidden
	}

	return s.userRepository.Update(ctx, update.ID, update.Name)
}
