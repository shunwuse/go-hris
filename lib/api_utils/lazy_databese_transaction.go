package api_utils

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/lib"
	"gorm.io/gorm"
)

type LazyDatabaseTransaction struct {
	logger lib.Logger
	db     *lib.Database
	trx    *gorm.DB

	once sync.Once
}

func NewLazyDatabaseTransaction(
	logger lib.Logger,
	db *lib.Database,
) LazyDatabaseTransaction {
	return LazyDatabaseTransaction{
		logger: logger,
		db:     db,
	}
}

func (l *LazyDatabaseTransaction) IsTransactionOpen() bool {
	return l.trx != nil
}

func (l *LazyDatabaseTransaction) beginTransaction() {
	l.once.Do(func() {
		// Begin database transaction
		trx := l.db.Begin()
		l.logger.Info("Begin database transaction")

		// Set transaction into struct
		l.trx = trx
	})
}

func (l *LazyDatabaseTransaction) getTransaction() *gorm.DB {
	if !l.IsTransactionOpen() {
		l.beginTransaction()
	}

	return l.trx
}

func SetLazyTransactionToContext(ctx *gin.Context, lazyTrx *LazyDatabaseTransaction) {
	ctx.Set(constants.DBTransaction, lazyTrx)
}

func GetTransactionFromContext(ctx *gin.Context) *gorm.DB {
	lazyTrx, exists := ctx.Get(constants.DBTransaction)
	if !exists {
		return nil
	}

	return lazyTrx.(*LazyDatabaseTransaction).getTransaction()
}
