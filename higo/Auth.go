package higo

import (
	"github.com/dengpju/higo-config/config"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/exception"
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
	// 空配置
	if nil == config.All() {
		return false
	}
	// 判断是否不需要鉴权
	if nil != config.Get("env.auth.NotAuth") {
		_, ok := config.Get("env.auth.NotAuth").(config.Configure)[flag]
		return ok
	}
	return false
}
