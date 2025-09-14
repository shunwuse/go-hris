package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/lib/api_utils"
)

type DBTrxMiddleware struct {
	logger lib.Logger
	db     lib.Database
}

func NewDBTrxMiddleware(
	logger lib.Logger,
	db lib.Database,
) DBTrxMiddleware {
	return DBTrxMiddleware{
		logger: logger,
		db:     db,
	}
}

func (m DBTrxMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Create lazy database transaction
		lazyTrx := api_utils.NewLazyDatabaseTransaction(m.logger, &m.db)

		// Set lazy transaction to context
		api_utils.SetLazyTransactionToContext(ctx, &lazyTrx)

		// Call next middleware
		ctx.Next()

		if !lazyTrx.IsTransactionOpen() {
			return
		}

		// Get transaction from context
		trx := api_utils.GetTransactionFromContext(ctx)

		// Check if it not specified http status code, then rollback transaction
		if ctx.Writer.Status() >= 400 {
			// Rollback transaction
			m.logger.Info("Rollback database transaction")
			trx.Rollback()
			return
		}

		// Commit transaction
		m.logger.Info("Commit database transaction")
		trx.Commit()
	}
}

func (m DBTrxMiddleware) Setup(router *gin.Engine) {
	m.logger.Info("Setting up database transaction middleware")

	router.Use(m.Handler())
}
