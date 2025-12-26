package infra

import (
	"context"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jmoiron/sqlx"
	"github.com/shunwuse/go-hris/ent/entgen"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	client *entgen.Client
	rawDB  *sqlx.DB
}

var globalDatabase *Database

func GetDatabase() *Database {
	if globalDatabase == nil {
		db := newDatabase(GetConfig(), GetLogger())
		globalDatabase = &db
	}

	return globalDatabase
}

func newDatabase(config Config, logger *Logger) Database {
	db, err := sqlx.Open(dialect.SQLite, config.SqliteDBPath)
	if err != nil {
		logger.Fatal("failed to open database", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	drv := entsql.OpenDB(dialect.SQLite, db.DB)
	client := entgen.NewClient(entgen.Driver(drv))

	logger.Info("database connected successfully")

	return Database{
		client: client,
		rawDB:  db,
	}
}

// GetClient returns the Ent client.
func (d *Database) GetClient(ctx context.Context) *entgen.Client {
	return d.client
}

// GetRawDB returns the raw sqlx.DB instance.
func (d *Database) GetRawDB(ctx context.Context) *sqlx.DB {
	return d.rawDB
}
