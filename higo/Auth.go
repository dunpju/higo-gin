package higo

import (
	"github.com/dengpju/higo-throw/throw"
)

// 是否空标记
func IsEmptyFlag(router *Router)  {
	if router.Flag() == "" {
		throw.Throw(router.RelativePath() + "未设置标记",0)
	}
}

// 是否不用鉴权
func IsNotAuth(flag string) bool {
	if "" == flag {
		return false
	}
	// 空配置
	if nil == ConfigAll() {
		return false
	}
	// 判断是否不需要鉴权
	if nil != Config("NotAuth") {
		_, ok := Config("NotAuth")[flag]
		return ok
	}
	return false
}
