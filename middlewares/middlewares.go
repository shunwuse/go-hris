package middlewares

import "github.com/gin-gonic/gin"

type CommonMiddlewares []ICommonMiddleware

type ICommonMiddleware interface {
	Setup(router *gin.Engine)
}

func NewCommonMiddlewares(
	dbTrxMiddleware DBTrxMiddleware,
) CommonMiddlewares {
	return CommonMiddlewares{
		dbTrxMiddleware,
	}
}

func (m CommonMiddlewares) Setup(router *gin.Engine) {
	for _, middleware := range m {
		middleware.Setup(router)
	}
}
