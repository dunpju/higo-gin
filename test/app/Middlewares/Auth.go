package Middlewares

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-gin/test/app/Consts"
	"github.com/dunpju/higo-throw/exception"
	"github.com/gin-gonic/gin"
)

// Auth 鉴权
type Auth struct{}

// NewAuth 构造函数
func NewAuth() *Auth {
	return &Auth{}
}

func (this *Auth) Middle(hg *higo.Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		if route, ok := hg.GetRoute(cxt.Request.Method, cxt.Request.URL.Path); ok {
			// TODO::非静态页面需要鉴权
			if !higo.IsNotAuth(route.Flag()) && !route.IsStatic() && route.IsAuth() {
				if "" == cxt.GetHeader("X-Token") {
					exception.Throw(exception.Message(Consts.InvalidToken.Message()), exception.Code(int(Consts.InvalidToken)))
				}
			}
			cxt.Next()
		} else {
			exception.Throw(exception.Message(Consts.InvalidToken.Message()), exception.Code(int(Consts.InvalidToken)))
		}
	}
}
