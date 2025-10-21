package services

import (
	"context"
	"errors"
	"log/slog"
	"slices"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/user"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
)

type userService struct {
	userRepository           repositories.UserRepository
	roleRepository           repositories.RoleRepository
	userRoleRepository       repositories.UserRoleRepository
	rolePermissionRepository repositories.RolePermissionRepository
}

func NewUserService(
	userRepository repositories.UserRepository,
	roleRepository repositories.RoleRepository,
	userRoleRepository repositories.UserRoleRepository,
	rolePermissionRepository repositories.RolePermissionRepository,
) service.UserService {
	return userService{
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
		slog.Error("Error getting users", "error", err)
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
		slog.Error("Error creating user", "error", err)
		return err
	}

	_, err = s.userRepository.Client.Password.
		Create().
		SetHash(user.Password.Hash).
		SetOwner(u).
		Save(ctx)
	if err != nil {
		slog.Error("Error creating password", "error", err)
		return err
	}

	roleModel := s.roleRepository.GetRoleByName(ctx, role.String())
	if roleModel == nil {
		// slog.Info("Role not found, creating role", "role", role)

		// roleCreate := &domains.RoleCreate{
		// 	Name: constants.Staff.String(),
		// }

		// if err := s.roleRepository.AddRole(ctx, roleCreate); err != nil {
		// 	slog.Error("add role error", "error", err)
		// 	return err
		// }

		slog.Error("Role not found", "role", role)
		return errors.New("role not found")
	}

	// Create user role
	_, err = s.userRepository.Client.UserRole.
		Create().
		SetUserID(u.ID).
		SetRoleID(roleModel.ID).
		Save(ctx)
	if err != nil {
		slog.Error("creating user role error", "error", err)
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
		slog.Error("Error getting user by username", "error", err)
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
		slog.Error("Error updating user", "error", err)
		return err
	}

	return nil
}
