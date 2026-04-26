package repositories

import (
	"context"
	"math"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/role"
	"github.com/shunwuse/go-hris/ent/entgen/user"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/repository"
	"go.uber.org/zap"
)

type UserRepository struct {
	logger *logger.Logger
	*database.Database
	cache *cache.Cache
}

func NewUserRepository(
	log *logger.Logger,
	db *database.Database,
	c *cache.Cache,
) repository.UserRepository {
	return &UserRepository{
		logger:   log,
		Database: db,
		cache:    c,
	}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]domains.User, error) {
	users, err := r.GetClient(ctx).User.Query().
		All(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find all users", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return mapUsers(users), nil
}

func (r *UserRepository) FindAllWithOffset(ctx context.Context, query domains.OffsetQuery, filter domains.UserFilter) (*domains.OffsetResult[domains.User], error) {
	dbQuery := r.GetClient(ctx).User.Query()

	if filter.Name != "" {
		dbQuery = dbQuery.Where(
			user.Or(
				user.UsernameHasPrefix(filter.Name),
				user.NameContainsFold(filter.Name),
			),
		)
	}

	if filter.Role != "" {
		dbQuery = dbQuery.Where(
			user.HasRolesWith(
				role.NameEQ(filter.Role),
			),
		)
	}

	total, err := dbQuery.Count(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to count users", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	users, err := dbQuery.
		Limit(query.PerPage).
		Offset(query.Offset()).
		Order(entgen.Desc(user.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find users with offset", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return &domains.OffsetResult[domains.User]{
		Items:      mapUsers(users),
		TotalCount: total,
		TotalPage:  int(math.Ceil(float64(total) / float64(query.PerPage))),
	}, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domains.User, error) {
	u, err := r.GetClient(ctx).User.Query().
		Where(user.UsernameEQ(username)).
		Only(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find user by username", zap.Error(err), zap.String("username", username))
		if entgen.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrDatabaseError
	}

	return mapUser(u), nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*domains.UserWithRoles, error) {
	return cache.Fetch(ctx, r.cache, constants.GetUserKey(id), constants.CacheTTLUser, func() (*domains.UserWithRoles, error) {
		u, err := r.GetClient(ctx).User.Query().
			WithPassword().
			WithRoles().
			Where(user.IDEQ(id)).
			Only(ctx)
		if err != nil {
			r.logger.WithContext(ctx).Error("failed to find user by id", zap.Error(err), zap.Uint("user_id", id))
			if entgen.IsNotFound(err) {
				return nil, errors.ErrNotFound
			}
			return nil, errors.ErrDatabaseError
		}

		return mapUserWithRoles(u), nil
	})
}

func (r *UserRepository) Create(ctx context.Context, username string, name string) (*domains.User, error) {
	u, err := r.GetClient(ctx).User.Create().
		SetUsername(username).
		SetName(name).
		Save(ctx)
	if err != nil {
		if entgen.IsConstraintError(err) {
			return nil, errors.ErrAlreadyExists
		}
		r.logger.WithContext(ctx).Error("failed to create user", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return mapUser(u), nil
}

func (r *UserRepository) Update(ctx context.Context, id uint, name string) error {
	if err := r.GetClient(ctx).User.Update().
		Where(user.IDEQ(id)).
		SetName(name).
		Exec(ctx); err != nil {
		r.logger.WithContext(ctx).Error("failed to update user", zap.Error(err), zap.Uint("user_id", id))
		return errors.ErrDatabaseError
	}

	r.invalidateUserCache(ctx, id)

	return nil
}

func (r *UserRepository) CreatePassword(ctx context.Context, hash string, ownerID uint) error {
	_, err := r.GetClient(ctx).Password.Create().
		SetHash(hash).
		SetOwnerID(ownerID).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create password", zap.Error(err))
		return errors.ErrDatabaseError
	}

	return nil
}

func (r *UserRepository) AssignRole(ctx context.Context, userID uint, roleID uint) error {
	_, err := r.GetClient(ctx).UserRole.Create().
		SetUserID(userID).
		SetRoleID(roleID).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create user role association", zap.Error(err))
		return errors.ErrDatabaseError
	}

	r.invalidateUserCache(ctx, userID)

	return nil
}

func (r *UserRepository) invalidateUserCache(ctx context.Context, userID uint) {
	r.cache.Client.Del(ctx, constants.GetUserKey(userID))
	r.cache.Client.Del(ctx, constants.GetUserPermissionsKey(userID))
}

func mapUsers(users []*entgen.User) []domains.User {
	result := make([]domains.User, len(users))
	for idx, user := range users {
		result[idx] = *mapUser(user)
	}

	return result
}

func mapUser(user *entgen.User) *domains.User {
	if user == nil {
		return nil
	}

	return &domains.User{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func mapUserWithRoles(user *entgen.User) *domains.UserWithRoles {
	if user == nil {
		return nil
	}

	roles := make([]constants.Role, len(user.Edges.Roles))
	for idx, role := range user.Edges.Roles {
		roles[idx] = role.Name
	}

	passwordHash := ""
	if user.Edges.Password != nil {
		passwordHash = user.Edges.Password.Hash
	}

	return &domains.UserWithRoles{
		ID:           user.ID,
		Username:     user.Username,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		PasswordHash: passwordHash,
		Roles:        roles,
	}
}
