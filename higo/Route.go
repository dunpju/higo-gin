package higo

import (
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/exception"
)

var RouterContainer RouterCollect

type RouterCollect map[string]*router.Route

// 添加路由容器
func (this RouterCollect) Add(relativePath string, route *router.Route) *RouterCollect {
	this[relativePath] = route
	return &this
}

// 所有路由
func (this RouterCollect) All() RouterCollect {
	return this
}

// 获取路由
func (this RouterCollect) Get(relativePath string) *router.Route {
	route, ok := this[relativePath]
	if !ok {
		exception.Throw(exception.Message(relativePath+"未定义路由"), exception.Code(0))
	}
	return route
}
