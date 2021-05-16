package higo

import (
	"github.com/gin-gonic/gin"
)

// 中间件接口(实现该接口都认为是中间件)
type IMiddleware interface {
	Middle(hg *Higo) gin.HandlerFunc
}

type Middleware struct{}

func (this *Middleware) Middle(hg *Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		cxt.Next()
	}
}
