package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/role"
	"github.com/shunwuse/go-hris/ent/entgen/user"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/pkg/config"
	"github.com/shunwuse/go-hris/internal/pkg/cryptox"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"go.uber.org/zap"
)

const (
	defaultUsers     = 20
	defaultApprovals = 50
	defaultSeed      = int64(20240624)
	demoUsername     = "demo"
)

func main() {
	users := flag.Int("users", defaultUsers, "number of demo users to add")
	approvals := flag.Int("approvals", defaultApprovals, "number of demo approvals to add")
	seed := flag.Int64("seed", defaultSeed, "seed for reproducible approval assignments")
	flag.Parse()

	if *users < 0 || *approvals < 0 {
		log.Fatal("users and approvals must not be negative")
	}

	manager := config.NewManager[Config](
		config.WithDotEnv("configs/default.env"),
		config.WithDotEnv("configs/.env"),
		config.WithEnv("APP_", "."),
	)
	if err := manager.Load(); err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log := logger.New(logger.WithConfig(logger.DefaultConfig()), logger.WithConsole(true))
	db := database.New(&manager.Config().Database, log)
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("failed to close database", zap.Error(err))
		}
	}()

	if err := run(context.Background(), db.GetClient(context.Background()), *users, *approvals, *seed); err != nil {
		log.Fatal("seed failed", zap.Error(err))
	}

	log.Info("seed completed", zap.Int("users", *users), zap.Int("approvals", *approvals))
}

func run(ctx context.Context, client *entgen.Client, usersToAdd, approvalsToAdd int, seed int64) error {
	roleIDs, err := findRoles(ctx, client)
	if err != nil {
		return err
	}

	if err := seedUsers(ctx, client, roleIDs, usersToAdd, seed); err != nil {
		return err
	}

	demoUsers, err := client.User.Query().
		Where(user.UsernameHasPrefix(demoUsername + ".")).
		WithRoles().
		Order(entgen.Asc(user.FieldID)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query demo users: %w", err)
	}

	admin, err := client.User.Query().
		Where(user.UsernameEQ("admin")).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("query admin user: %w", err)
	}

	return seedApprovals(ctx, client, demoUsers, admin.ID, approvalsToAdd, seed)
}

type roleIDs struct {
	manager uint
	staff   uint
}

func findRoles(ctx context.Context, client *entgen.Client) (roleIDs, error) {
	manager, err := client.Role.Query().Where(role.NameEQ(constants.Manager)).Only(ctx)
	if err != nil {
		return roleIDs{}, fmt.Errorf("query manager role: %w", err)
	}
	staff, err := client.Role.Query().Where(role.NameEQ(constants.Staff)).Only(ctx)
	if err != nil {
		return roleIDs{}, fmt.Errorf("query staff role: %w", err)
	}

	return roleIDs{manager: manager.ID, staff: staff.ID}, nil
}

func seedUsers(ctx context.Context, client *entgen.Client, roles roleIDs, count int, seed int64) error {
	managerCount := count * 3 / 10
	if count > 0 && managerCount == 0 {
		managerCount = 1
	}
	staffCount := count - managerCount
	passwordHash, err := cryptox.HashPassword("password")
	if err != nil {
		return fmt.Errorf("hash demo password: %w", err)
	}

	managerStart, err := nextDemoUserNumber(ctx, client, "demo.manager.")
	if err != nil {
		return err
	}
	staffStart, err := nextDemoUserNumber(ctx, client, "demo.staff.")
	if err != nil {
		return err
	}

	type userSeed struct {
		username string
		name     string
		roleID   uint
	}

	users := make([]userSeed, 0, count)
	for i := 0; i < managerCount; i++ {
		number := managerStart + i
		users = append(users, userSeed{
			username: fmt.Sprintf("demo.manager.%03d", number),
			name:     fmt.Sprintf("Demo Manager %03d", number),
			roleID:   roles.manager,
		})
	}
	for i := 0; i < staffCount; i++ {
		number := staffStart + i
		users = append(users, userSeed{
			username: fmt.Sprintf("demo.staff.%03d", number),
			name:     fmt.Sprintf("Demo Staff %03d", number),
			roleID:   roles.staff,
		})
	}

	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})
	for _, demoUser := range users {
		if err := createUser(ctx, client, demoUser.username, demoUser.name, demoUser.roleID, passwordHash); err != nil {
			return err
		}
	}

	return nil
}

func nextDemoUserNumber(ctx context.Context, client *entgen.Client, prefix string) (int, error) {
	users, err := client.User.Query().Where(user.UsernameHasPrefix(prefix)).All(ctx)
	if err != nil {
		return 0, fmt.Errorf("query demo users with prefix %q: %w", prefix, err)
	}

	maxNumber := 0
	for _, current := range users {
		var number int
		if _, err := fmt.Sscanf(current.Username, prefix+"%d", &number); err == nil && number > maxNumber {
			maxNumber = number
		}
	}

	return maxNumber + 1, nil
}

func createUser(ctx context.Context, client *entgen.Client, username, name string, roleID uint, passwordHash string) error {
	created, err := client.User.Create().
		SetUsername(username).
		SetName(name).
		AddRoleIDs(roleID).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("create user %q: %w", username, err)
	}
	if _, err := client.Password.Create().
		SetHash(passwordHash).
		SetOwnerID(created.ID).
		Save(ctx); err != nil {
		return fmt.Errorf("create password for %q: %w", username, err)
	}

	return nil
}

func seedApprovals(ctx context.Context, client *entgen.Client, demoUsers []*entgen.User, adminID uint, count int, seed int64) error {
	if len(demoUsers) == 0 && count > 0 {
		return fmt.Errorf("no demo users available for approvals")
	}

	managers, err := client.User.Query().
		Where(user.UsernameHasPrefix("demo.manager.")).
		Order(entgen.Asc(user.FieldID)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query demo managers: %w", err)
	}
	approvers := append([]*entgen.User{{ID: adminID}}, managers...)
	rng := rand.New(rand.NewSource(seed))
	statuses := make([]constants.ApprovalStatus, 0, count)
	statuses = append(statuses, repeatedStatus(constants.ApprovalStatusPending, count-count*15/100)...)
	statuses = append(statuses, repeatedStatus(constants.ApprovalStatusApproved, count*10/100)...)
	statuses = append(statuses, repeatedStatus(constants.ApprovalStatusRejected, count*5/100)...)
	rng.Shuffle(len(statuses), func(i, j int) {
		statuses[i], statuses[j] = statuses[j], statuses[i]
	})

	for index := 1; index <= count; index++ {
		creator := demoUsers[rng.Intn(len(demoUsers))]
		status := statuses[index-1]
		var approverID uint
		if status != constants.ApprovalStatusPending {
			approver := approvers[rng.Intn(len(approvers))]
			for offset := 0; approver.ID == creator.ID && offset < len(approvers); offset++ {
				approver = approvers[(index+offset)%len(approvers)]
			}
			approverID = approver.ID
		}

		builder := client.Approval.Create().
			SetStatus(status).
			SetCreatorID(creator.ID)
		if approverID != 0 {
			builder.SetApproverID(approverID)
		}
		if _, err := builder.Save(ctx); err != nil {
			return fmt.Errorf("create approval %d: %w", index, err)
		}
	}

	return nil
}

func repeatedStatus(status constants.ApprovalStatus, count int) []constants.ApprovalStatus {
	statuses := make([]constants.ApprovalStatus, count)
	for index := range statuses {
		statuses[index] = status
	}

	return statuses
}
