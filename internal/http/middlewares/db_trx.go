package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shunwuse/go-hris/internal/http/api_utils"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type DBTrxMiddleware struct {
	logger *infra.Logger
	db     *infra.Database
}

func NewDBTrxMiddleware(
	logger *infra.Logger,
	db *infra.Database,
) *DBTrxMiddleware {
	return &DBTrxMiddleware{
		logger: logger,
		db:     db,
	}
}

func (m *DBTrxMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create lazy database transaction.
			lazyTrx := api_utils.NewLazyDatabaseTransaction(m.logger, m.db)

			// Set lazy transaction to context.
			ctx := api_utils.SetLazyTransactionToContext(r.Context(), &lazyTrx)

			// Wrap response writer to capture status code (using Chi's official wrapper).
			writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Call next handler.
			next.ServeHTTP(writer, r.WithContext(ctx))

			if !lazyTrx.IsTransactionOpen() {
				return
			}

			// Get transaction from context.
			trx := api_utils.GetTransactionFromContext(ctx)

			// Check status code to decide commit or rollback.
			if writer.Status() >= 400 {
				// Rollback transaction on error.
				m.logger.WithContext(ctx).Info("rollback database transaction")
				if err := trx.Rollback(); err != nil {
					m.logger.WithContext(ctx).Error("failed to rollback transaction", zap.Error(err))
				}
				return
			}

			// Commit transaction.
			m.logger.WithContext(ctx).Info("commit database transaction")
			if err := trx.Commit(); err != nil {
				m.logger.WithContext(ctx).Error("failed to commit transaction", zap.Error(err))
			}
		})
	}
}

func (m *DBTrxMiddleware) Setup(router chi.Router) {
	m.logger.Info("setting up database transaction middleware")

	router.Use(m.Handler())
}
