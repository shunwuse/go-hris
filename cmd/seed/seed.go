package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/role"
	"github.com/shunwuse/go-hris/ent/entgen/user"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/pkg/cryptox"
)

func seedBatch(ctx context.Context, client *entgen.Client, users, approvals int, seed int64, password string) error {
	rng := rand.New(rand.NewPCG(uint64(seed), 0))

	if users > 0 {
		if err := addUsers(ctx, client, rng, users, password); err != nil {
			return err
		}
	}

	if approvals > 0 {
		return addApprovals(ctx, client, rng, approvals)
	}

	return nil
}

func demoPrefix(name constants.Role) string {
	return "demo." + name.String() + "."
}

func demoNumber(username, prefix string) int {
	suffix, ok := strings.CutPrefix(username, prefix)
	if !ok {
		return 0
	}

	number, err := strconv.Atoi(suffix)
	if err != nil || number <= 0 || fmt.Sprintf("%03d", number) != suffix {
		return 0
	}

	return number
}

func userRoles(count int) []constants.Role {
	roles := make([]constants.Role, count)
	managers := 0

	if count >= 2 {
		managers = max(1, count/3)
	}

	for i := range roles {
		roles[i] = constants.Staff
		if i < managers {
			roles[i] = constants.Manager
		}
	}

	return roles
}

func addUsers(ctx context.Context, client *entgen.Client, rng *rand.Rand, count int, password string) error {
	type allocation struct {
		roleID uint
		last   int
	}

	roles := userRoles(count)
	allocations := make(map[constants.Role]allocation, 2)

	for _, name := range roles {
		if _, ok := allocations[name]; ok {
			continue
		}

		roleID, err := client.Role.Query().
			Where(role.NameEQ(name), role.DeletedAtIsNil()).OnlyID(ctx)
		if err != nil {
			return fmt.Errorf("look up %s role (have you applied migrations?): %w", name, err)
		}

		prefix := demoPrefix(name)
		// Include deleted users: their usernames are still reserved by PostgreSQL.

		names, err := client.User.Query().Where(user.UsernameHasPrefix(prefix)).Select(user.FieldUsername).Strings(ctx)
		if err != nil {
			return fmt.Errorf("look up %s usernames: %w", name, err)
		}

		last := 0
		for _, username := range names {
			last = max(last, demoNumber(username, prefix))
		}

		allocations[name] = allocation{roleID: roleID, last: last}
	}

	rng.Shuffle(len(roles), func(i, j int) { roles[i], roles[j] = roles[j], roles[i] })

	for _, name := range roles {
		next := allocations[name]

		if next.last == math.MaxInt {
			return fmt.Errorf("no demo username numbers left for %s", name)
		}

		next.last++
		allocations[name] = next

		username := fmt.Sprintf("%s%03d", demoPrefix(name), next.last)

		if err := ctx.Err(); err != nil {
			return err
		}

		hash, err := cryptox.HashPassword(password)
		if err != nil {
			return fmt.Errorf("hash demo password: %w", err)
		}

		created, err := client.User.Create().
			SetUsername(username).
			SetName(fmt.Sprintf("Demo %s %03d", name, next.last)).
			AddRoleIDs(next.roleID).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("create user %q: %w", username, err)
		}

		if _, err := client.Password.Create().SetOwnerID(created.ID).SetHash(hash).Save(ctx); err != nil {
			return fmt.Errorf("create password for %q: %w", username, err)
		}
	}

	return nil
}

func demoUserIDs(ctx context.Context, client *entgen.Client, name constants.Role) ([]uint, error) {
	prefix := demoPrefix(name)

	users, err := client.User.Query().
		Where(user.DeletedAtIsNil(), user.UsernameHasPrefix(prefix),
			user.HasRolesWith(role.NameEQ(name), role.DeletedAtIsNil())).
		Order(user.ByID()).
		Select(user.FieldID, user.FieldUsername).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("look up demo %s users: %w", name, err)
	}

	var ids []uint
	for _, demo := range users {
		if demoNumber(demo.Username, prefix) > 0 {
			ids = append(ids, demo.ID)
		}
	}

	return ids, nil
}

func approvalStatuses(count int) []constants.ApprovalStatus {
	approved, rejected := count/10, count/20

	statuses := make([]constants.ApprovalStatus, count)
	for i := range statuses {
		switch {
		case i < approved:
			statuses[i] = constants.ApprovalStatusApproved
		case i < approved+rejected:
			statuses[i] = constants.ApprovalStatusRejected
		default:
			statuses[i] = constants.ApprovalStatusPending
		}
	}

	return statuses
}

func addApprovals(ctx context.Context, client *entgen.Client, rng *rand.Rand, count int) error {
	creators, err := demoUserIDs(ctx, client, constants.Staff)
	if err != nil {
		return err
	}

	if len(creators) == 0 {
		return errors.New("no active demo staff available: add demo users before seeding approvals")
	}

	statuses := approvalStatuses(count)

	var approvers []uint
	if slices.ContainsFunc(statuses, func(status constants.ApprovalStatus) bool {
		return status != constants.ApprovalStatusPending
	}) {
		approvers, err = demoUserIDs(ctx, client, constants.Manager)
		if err != nil {
			return err
		}

		if len(approvers) == 0 {
			admin, err := client.User.Query().
				Where(user.UsernameEQ("admin"), user.DeletedAtIsNil(),
					user.HasRolesWith(role.NameEQ(constants.Admin), role.DeletedAtIsNil())).
				OnlyID(ctx)
			if err != nil {
				return fmt.Errorf("no demo managers; look up active administrator account: %w", err)
			}

			approvers = []uint{admin}
		}
	}

	rng.Shuffle(len(statuses), func(i, j int) { statuses[i], statuses[j] = statuses[j], statuses[i] })

	for i, status := range statuses {
		builder := client.Approval.Create().SetStatus(status).SetCreatorID(creators[rng.IntN(len(creators))])

		if status != constants.ApprovalStatusPending {
			builder.SetApproverID(approvers[rng.IntN(len(approvers))])
		}

		if _, err := builder.Save(ctx); err != nil {
			return fmt.Errorf("create demo approval %d: %w", i+1, err)
		}
	}

	return nil
}
