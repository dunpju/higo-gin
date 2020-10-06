package higo

import (
	"github.com/gin-gonic/gin"
)

// 中间件接口(实现该接口都认为是中间件)
type IMiddleware interface {
	Loader(hg *Higo) gin.HandlerFunc
}
