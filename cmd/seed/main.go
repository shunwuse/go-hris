package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/pkg/config"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
)

const (
	defaultUsers     = 20
	defaultApprovals = 50
	defaultSeed      = int64(20240624)
	defaultPassword  = "password"
)

type seedConfig struct {
	Environment string          `koanf:"env"`
	Database    database.Config `koanf:",squash"`
}

func main() {
	users := flag.Int("users", defaultUsers, "number of demo users to append")
	approvals := flag.Int("approvals", defaultApprovals, "number of demo approvals to append")
	seed := flag.Int64("seed", defaultSeed, "random seed for role order, statuses, and approval assignments")

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: go run ./cmd/seed [flags]")
		fmt.Fprintf(flag.CommandLine.Output(), "For disposable development databases only. Set SEED_PASSWORD to override the demo password (default: %q).\n", defaultPassword)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 0 {
		log.Fatal("seed: unexpected positional arguments")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := execute(ctx, *users, *approvals, *seed)
	stop()
	if err != nil {
		log.Fatalf("seed: %v", err)
	}

	log.Printf("seeded %d demo users and %d demo approvals", *users, *approvals)
}

func execute(ctx context.Context, users, approvals int, seed int64) (err error) {
	if users < 0 || approvals < 0 {
		return errors.New("-users and -approvals must not be negative")
	}

	if users == 0 && approvals == 0 {
		return nil
	}

	password := os.Getenv("SEED_PASSWORD")
	if password == "" {
		password = defaultPassword
	}

	if users > 0 && len(password) > 72 {
		return errors.New("SEED_PASSWORD must not exceed 72 bytes when adding demo users")
	}

	manager := config.NewManager[seedConfig](
		config.WithDotEnv("configs/default.env"),
		config.WithDotEnv("configs/.env"),
		config.WithEnv("APP_", "."),
	)

	if err := manager.Load(); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	cfg := manager.Config()
	if cfg.Environment != constants.EnvDevelopment {
		return errors.New("seeding is only allowed with env=development")
	}

	db := database.New(&cfg.Database, logger.New(logger.WithConsole(true)))
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close database: %w", closeErr))
		}
	}()

	return db.WithTx(ctx, func(ctx context.Context) error {
		return seedBatch(ctx, db.GetClient(ctx), users, approvals, seed, password)
	})
}
