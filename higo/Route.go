package higo

import (
	"github.com/dengpju/higo-throw/throw"
)

type Route struct {
	Method       string      // 请求方法 GET/POST/DELETE/PATCH/OPTIONS/HEAD
	RelativePath string      // 后端 api relativePath
	Handle       interface{} // 后端控制器函数
	Flag         string      // 后端控制器函数标记
	FrontPath    string      // 前端 path(前端菜单路由)
	IsStatic     bool        // 是否静态文件
	Desc         string      // 描述
}

func NewRoute() *Route {
	return &Route{}
}

var Router RouterCollect

type RouterCollect map[string]Route

func NewRouter() *RouterCollect {
	return &Router
}

// 添加路由容器
func (this RouterCollect) Add(relativePath string, route Route) *RouterCollect {
	this[relativePath] = route
	return &this
}

// 所有路由
func (this RouterCollect) All() map[string]Route {
	return this
}

// 获取路由
func (this RouterCollect) Get(relativePath string) Route {
	route, ok := this[relativePath]
	if !ok {
		throw.Throw(relativePath+"未定义路由", 0)
	}
	return route
}

func (this Route) Get(relativePath string, handle interface{}) Route {
	this.Method = "Get"
	this.RelativePath = relativePath
	this.Handle = handle
	this.Handle = handle
	return this
}
