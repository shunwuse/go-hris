package api_utils

import (
	"context"
	"sync"

	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/lib"
)

type LazyDatabaseTransaction struct {
	logger lib.Logger
	db     *lib.Database
	trx    *entgen.Tx

	once sync.Once
}

func NewLazyDatabaseTransaction(
	logger lib.Logger,
	client *lib.Database,
) LazyDatabaseTransaction {
	return LazyDatabaseTransaction{
		logger: logger,
		db:     client,
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
			l.logger.Errorf("Failed to begin transaction: %v", err)
		}

		l.logger.Info("Begin database transaction")

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
