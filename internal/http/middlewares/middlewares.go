package middlewares

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type CommonMiddlewares []ICommonMiddleware

type ICommonMiddleware interface {
	Setup(router chi.Router)
}

func NewCommonMiddlewares(
	dbTrxMiddleware DBTrxMiddleware,
) CommonMiddlewares {
	return CommonMiddlewares{
		dbTrxMiddleware,
	}
}

func (m CommonMiddlewares) Setup(router chi.Router) {
	// Built-in middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	for _, middleware := range m {
		middleware.Setup(router)
	}
}
