package higo

import (
	"gitee.com/dengpju/higo-configure/configure"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/throw"
)

// 是否空标记
func IsEmptyFlag(route *router.Route) {
	if route.Flag() == "" && !route.IsStatic() {
		throw.Throw(throw.Message(route.RelativePath()+"未设置标记"), throw.Code(0))
	}
}

// 是否不用鉴权
func IsNotAuth(flag string) bool {
	if "" == flag {
		return false
	}
	// 空配置
	if nil == configure.All() {
		return false
	}
	// 判断是否不需要鉴权
	if nil != configure.Config("NotAuth") {
		_, ok := configure.Config("NotAuth")[flag]
		return ok
	}
	return false
}
