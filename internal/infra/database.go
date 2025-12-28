package infra

import (
	"context"
	"fmt"
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

// WithTx provides a transactional context to the given work function.
func (d *Database) WithTx(ctx context.Context, work func(ctx context.Context) error) error {
	// 1. Begin transaction.
	tx, err := d.client.Tx(ctx)
	if err != nil {
		return err
	}

	// 2. Handle panic.
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 3. Execute work with transactional context.
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err := work(ctxWithTx); err != nil {
		// 4. Rollback on error.
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}

		return err
	}

	// 5. Commit on success.
	return tx.Commit()
}

type txKey struct{}
