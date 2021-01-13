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

type Router struct {
	method       string      // 请求方法 GET/POST/DELETE/PATCH/OPTIONS/HEAD
	relativePath string      // 后端 api relativePath
	Handle       interface{} // 后端控制器函数
	flag         string      // 后端控制器函数标记
	frontPath    string      // 前端 path(前端菜单路由)
	isStatic     bool        // 是否静态文件
	desc         string      // 描述
}

func Route(args ...*RouteAttribute) Router {
	router := &Router{}
	for _, attribute := range args {
		if attribute.name == ROUTE_METHOD {
			router.method = attribute.value.(string)
		} else if attribute.name == ROUTE_RELATIVE_PATH {
			router.relativePath = attribute.value.(string)
		} else if attribute.name == ROUTE_HANDLE {
			router.Handle = attribute.value
		} else if attribute.name == ROUTE_FLAG {
			router.flag = attribute.value.(string)
		} else if attribute.name == ROUTE_FRONTPATH {
			router.frontPath = attribute.value.(string)
		} else if attribute.name == ROUTE_DESC {
			router.desc = attribute.value.(string)
		}
	}
	return *router
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

var RouterContainer RouterCollect

type RouterCollect map[string]Router

// 添加路由容器
func (this RouterCollect) Add(relativePath string, router Router) *RouterCollect {
	this[relativePath] = router
	return &this
}

// 所有路由
func (this RouterCollect) All() RouterCollect {
	return this
}

// 获取路由
func (this RouterCollect) Get(relativePath string) Router {
	router, ok := this[relativePath]
	if !ok {
		throw.Throw(relativePath+"未定义路由", 0)
	}
	return router
}

func (this Router) Method() string {
	return this.method
}

func (this Router) RelativePath() string {
	return this.relativePath
}

func (this Router) Flag() string {
	return this.flag
}

func (this Router) FrontPath() string {
	return this.frontPath
}

func (this Router) IsStatic() bool {
	return this.isStatic
}

func (this Router) Desc() string {
	return this.desc
}

func (this Router) Get(relativePath string, handle interface{}) Router {
	this.method = "GET"
	this.relativePath = relativePath
	this.Handle = handle
	this.Handle = handle
	return this
}
