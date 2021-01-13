package higo

import (
	"github.com/dengpju/higo-throw/throw"
)

const (
	ROUTE_METHOD        = "method"
	ROUTE_RELATIVE_PATH = "relativePath"
	ROUTE_HANDLE        = "handle"
	ROUTE_FLAG          = "flag"
	ROUTE_FRONTPATH     = "frontPath"
	ROUTE_IS_STATIC     = "isStatic"
	ROUTE_DESC          = "desc"
)

type Route struct {
	method       string      // 请求方法 GET/POST/DELETE/PATCH/OPTIONS/HEAD
	relativePath string      // 后端 api relativePath
	Handle       interface{} // 后端控制器函数
	flag         string      // 后端控制器函数标记
	frontPath    string      // 前端 path(前端菜单路由)
	isStatic     bool        // 是否静态文件
	desc         string      // 描述
}

func NewRoute(args ...*RouteAttribute) Route {
	route := &Route{}
	for _, attribute := range args {
		if attribute.name == ROUTE_METHOD {
			route.method = attribute.value.(string)
		} else if attribute.name == ROUTE_RELATIVE_PATH {
			route.relativePath = attribute.value.(string)
		} else if attribute.name == ROUTE_HANDLE {
			route.Handle = attribute.value
		} else if attribute.name == ROUTE_FLAG {
			route.flag = attribute.value.(string)
		} else if attribute.name == ROUTE_FRONTPATH {
			route.frontPath = attribute.value.(string)
		} else if attribute.name == ROUTE_DESC {
			route.desc = attribute.value.(string)
		}
	}
	return *route
}

type RouteAttribute struct {
	name  string
	value interface{}
}

func NewRouteAttribute(name string, value interface{}) *RouteAttribute {
	return &RouteAttribute{name: name, value: value}
}

func Method(value string) *RouteAttribute {
	return NewRouteAttribute(ROUTE_METHOD, value)
}

func RelativePath(value string) *RouteAttribute {
	return NewRouteAttribute(ROUTE_RELATIVE_PATH, value)
}

func Handle(value interface{}) *RouteAttribute {
	return NewRouteAttribute(ROUTE_HANDLE, value)
}

func Flag(value string) *RouteAttribute {
	return NewRouteAttribute(ROUTE_FLAG, value)
}

func FrontPath(value string) *RouteAttribute {
	return NewRouteAttribute(ROUTE_FRONTPATH, value)
}

func IsStatic(value bool) *RouteAttribute {
	return NewRouteAttribute(ROUTE_IS_STATIC, value)
}

func Desc(value string) *RouteAttribute {
	return NewRouteAttribute(ROUTE_DESC, value)
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

func (this Route) Method() string {
	return this.method
}

func (this Route) RelativePath() string {
	return this.relativePath
}

func (this Route) Flag() string {
	return this.flag
}

func (this Route) FrontPath() string {
	return this.frontPath
}

func (this Route) IsStatic() bool {
	return this.isStatic
}

func (this Route) Desc() string {
	return this.desc
}

func (this Route) Get(relativePath string, handle interface{}) Route {
	this.method = "GET"
	this.relativePath = relativePath
	this.Handle = handle
	this.Handle = handle
	return this
}
