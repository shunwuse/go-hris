package database

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jmoiron/sqlx"
	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra/config"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

type Database struct {
	client *entgen.Client
	rawDB  *sqlx.DB
}

var (
	instance *Database
)

// DB returns the database client.
func DB() *Database {
	if instance == nil {
		cfg, _ := config.Load()
		log := logger.L()
		if log == nil {
			log = logger.New(
				logger.WithConfig(logger.Config{
					Level:      cfg.Log.Level,
					FilePath:   cfg.Log.FilePath,
					MaxSize:    cfg.Log.MaxSize,
					MaxBackups: cfg.Log.MaxBackups,
					MaxAge:     cfg.Log.MaxAge,
					Compress:   cfg.Log.Compress,
				}),
				logger.WithConsole(cfg.Service.Environment == constants.EnvDevelopment),
			)
		}
		db := connect(cfg, log)
		instance = &db
	}

	return instance
}

func connect(cfg *config.Config, log *logger.Logger) Database {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode,
	)

	db, err := sqlx.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatal("failed to open database", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("failed to ping database", zap.Error(err))
	}

	drv := entsql.OpenDB(dialect.Postgres, db.DB)
	client := entgen.NewClient(entgen.Driver(drv))

	log.Info("database connected successfully")

	return Database{
		client: client,
		rawDB:  db,
	}
}

// Close closes the database connection.
func (d *Database) Close() error {
	if d.client != nil {
		if err := d.client.Close(); err != nil {
			return err
		}
	}
	return d.rawDB.Close()
}

// GetClient returns the transactional or default database client from context.
func (d *Database) GetClient(ctx context.Context) *entgen.Client {
	tx, ok := ctx.Value(txKey{}).(*entgen.Tx)
	if ok {
		return tx.Client()
	}

	return d.client
}

// GetRawDB returns the raw sqlx.DB instance.
func (d *Database) GetRawDB(ctx context.Context) *sqlx.DB {
	return d.rawDB
}
