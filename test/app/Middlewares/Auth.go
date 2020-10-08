package Middlewares

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Consts"
	"github.com/gin-gonic/gin"
)

// 鉴权
type Auth struct {}

// 构造函数
func NewAuth() *Auth {
	return &Auth{}
}

func (this Auth) Loader(hg *higo.Higo) gin.HandlerFunc {
	return func(c *gin.Context) {
		if route, ok := hg.GetRoute(c.Request.URL.Path); ok {
			if !higo.IsNotAuth(route.Flag) {
				if "" == c.GetHeader("X-Token") {
					higo.Throw(higo.Const(Consts.INVALID_TOKEN).Msg, higo.Const(Consts.INVALID_TOKEN).Code)
				}
			}
			c.Next()
		}else {
			higo.Throw(higo.Const(Consts.INVALID_API).Msg, higo.Const(Consts.INVALID_API).Code)
		}
	}
}
