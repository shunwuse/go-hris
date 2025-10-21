package api_utils

import (
	"context"
	"log/slog"
	"sync"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra"
)

type LazyDatabaseTransaction struct {
	db  *infra.Database
	trx *entgen.Tx

	once sync.Once
}

func NewLazyDatabaseTransaction(
	client *infra.Database,
) LazyDatabaseTransaction {
	return LazyDatabaseTransaction{
		db: client,
	}
}

func (l *LazyDatabaseTransaction) IsTransactionOpen() bool {
	return l.trx != nil
}

func (l *LazyDatabaseTransaction) beginTransaction() {
	l.once.Do(func() {
		// Begin database transaction
		trx, err := l.db.Client.Tx(context.Background())
		if err != nil {
			slog.Error("Failed to begin transaction", "error", err)
		}

		slog.Info("Begin database transaction")

		// Set transaction into struct
		l.trx = trx
	})
}

func (l *LazyDatabaseTransaction) getTransaction() *entgen.Tx {
	if !l.IsTransactionOpen() {
		l.beginTransaction()
	}

	return l.trx
}

func SetLazyTransactionToContext(ctx context.Context, lazyTrx *LazyDatabaseTransaction) context.Context {
	return context.WithValue(ctx, constants.DBTransaction, lazyTrx)
}

func GetTransactionFromContext(ctx context.Context) *entgen.Tx {
	lazyTrx, ok := ctx.Value(constants.DBTransaction).(*LazyDatabaseTransaction)
	if !ok || lazyTrx == nil {
		return nil
	}

	return lazyTrx.getTransaction()
}
