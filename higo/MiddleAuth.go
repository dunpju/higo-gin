package higo

import (
	"github.com/dengpju/higo-config/config"
	"github.com/dengpju/higo-gin/test/app/Consts"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/exception"
	"github.com/gin-gonic/gin"
)

// 是否空标记
func IsEmptyFlag(route *router.Route) {
	if route.Flag() == "" && !route.IsStatic() {
		exception.Throw(exception.Message(route.RelativePath()+"未设置标记"), exception.Code(0))
	}
}

// 是否不用鉴权
func IsNotAuth(flag string) bool {
	if "" == flag {
		return false
	}
	// 判断是否不需要鉴权
	return config.Auth("NotAuth").(*config.Configure).Exist(flag)
}

// 鉴权
type Auth struct{}

// 构造函数
func NewAuth() *Auth {
	return &Auth{}
}

func (this *Auth) Middle(hg *Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		MiddleAuthFunc(hg)(cxt)
		cxt.Next()
	}
}

func middleAuthFunc(hg *Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		if route, ok := hg.GetRoute(cxt.Request.URL.Path); ok {
			if ! IsNotAuth(route.Flag()) && !route.IsStatic() {
				if "" == cxt.GetHeader("X-Token") {
					exception.Throw(exception.Message(Consts.InvalidToken.Message()), exception.Code(int(Consts.InvalidToken)))
				}
			}
		} else {
			exception.Throw(exception.Message(Consts.InvalidApi.Message()), exception.Code(int(Consts.InvalidApi)))
		}
	}
}