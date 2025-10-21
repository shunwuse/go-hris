package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shunwuse/go-hris/internal/http/api_utils"
	"github.com/shunwuse/go-hris/internal/infra"
)

type DBTrxMiddleware struct {
	db infra.Database
}

func NewDBTrxMiddleware(
	db infra.Database,
) DBTrxMiddleware {
	return DBTrxMiddleware{
		db: db,
	}
}

func (m DBTrxMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create lazy database transaction
			lazyTrx := api_utils.NewLazyDatabaseTransaction(&m.db)

			// Set lazy transaction to context
			ctx := api_utils.SetLazyTransactionToContext(r.Context(), &lazyTrx)

			// Wrap response writer to capture status code (using Chi's official wrapper)
			writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Call next handler
			next.ServeHTTP(writer, r.WithContext(ctx))

			if !lazyTrx.IsTransactionOpen() {
				return
			}

			// Get transaction from context
			trx := api_utils.GetTransactionFromContext(ctx)

			// Check status code to decide commit or rollback
			if writer.Status() >= 400 {
				// Rollback transaction on error
				slog.Info("Rollback database transaction")
				trx.Rollback()
				return
			}

			// Commit transaction
			slog.Info("Commit database transaction")
			trx.Commit()
		})
	}
}

func (m DBTrxMiddleware) Setup(router chi.Router) {
	slog.Info("Setting up database transaction middleware")

	router.Use(m.Handler())
}
