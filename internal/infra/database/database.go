package database

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jmoiron/sqlx"
	"github.com/shunwuse/go-hris/ent/entgen"
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
		db := connect(config.Get(), logger.L())
		instance = &db
	}

	return instance
}

func connect(cfg config.Config, log *logger.Logger) Database {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
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
