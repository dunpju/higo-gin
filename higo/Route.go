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
	handle       interface{} // 后端控制器函数
	flag         string      // 后端控制器函数标记
	frontPath    string      // 前端 path(前端菜单路由)
	isStatic     bool        // 是否静态文件
	desc         string      // 描述
}

func NewRoute(args ...RouteAttribute) *Route {
	route := &Route{}
	for _, attribute := range args {
		
		route.{attribute.name} = attribute.value
	}
	return route
}

type RouteAttribute struct {
	name  string
	value interface{}
}

func NewRouteAttribute(name string, value interface{}) *RouteAttribute {
	return &RouteAttribute{name: name, value: value}
}

func RouteMethod(value string) *RouteAttribute {
	return NewRouteAttribute(ROUTE_METHOD, value)
}

func RouteRelativePath(value string) *RouteAttribute {
	return NewRouteAttribute(ROUTE_RELATIVE_PATH, value)
}

func RouteHandle(value interface{}) *RouteAttribute {
	return NewRouteAttribute(ROUTE_HANDLE, value)
}

func RouteFlag(value string) *RouteAttribute {
	return NewRouteAttribute(ROUTE_FLAG, value)
}

func RouteFrontPath(value string) *RouteAttribute {
	return NewRouteAttribute(ROUTE_FRONTPATH, value)
}

func RouteIsStatic(value bool) *RouteAttribute {
	return NewRouteAttribute(ROUTE_IS_STATIC, value)
}

func RouteDesc(value string) *RouteAttribute {
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

func (this Route) Handle() interface{} {
	return this.handle
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
	this.handle = handle
	this.handle = handle
	return this
}
