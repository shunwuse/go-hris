package database

import (
	"context"
	"fmt"
)

type Transactor interface {
	WithTx(ctx context.Context, work func(ctx context.Context) error) error
}

// NewTransactor returns the database instance as a Transactor.
func NewTransactor() Transactor {
	return DB()
}

type txKey struct{}

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
