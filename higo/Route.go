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

type GroupRouter interface {
	This() interface{}
}

type Router struct {
	method       string      // 请求方法 GET/POST/DELETE/PATCH/OPTIONS/HEAD
	relativePath string      // 后端 api relativePath
	handle       interface{} // 后端控制器函数
	flag         string      // 后端控制器函数标记
	frontPath    string      // 前端 path(前端菜单路由)
	isStatic     bool        // 是否静态文件
	desc         string      // 描述
}

func (this *Router) This() interface{} {
	return this
}

func Route(args ...*RouteAttribute) *Router {
	router := &Router{}
	for _, attribute := range args {
		if attribute.Name == ROUTE_METHOD {
			router.method = attribute.Value.(string)
		} else if attribute.Name == ROUTE_RELATIVE_PATH {
			router.relativePath = attribute.Value.(string)
		} else if attribute.Name == ROUTE_HANDLE {
			router.handle = attribute.Value
		} else if attribute.Name == ROUTE_FLAG {
			router.flag = attribute.Value.(string)
		} else if attribute.Name == ROUTE_FRONTPATH {
			router.frontPath = attribute.Value.(string)
		} else if attribute.Name == ROUTE_DESC {
			router.desc = attribute.Value.(string)
		} else if attribute.Name == ROUTE_IS_STATIC {
			router.isStatic = attribute.Value.(bool)
		}
	}
	return router
}

type RouteAttributes []*RouteAttribute

func (this RouteAttributes) Find(name string) interface{} {
	for _, p := range this {
		if p.Name == name {
			return p.Value
		}
	}
	return nil
}

type RouteAttribute struct {
	Name  string
	Value interface{}
}

func NewRouteAttribute(name string, value interface{}) *RouteAttribute {
	return &RouteAttribute{Name: name, Value: value}
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

type RouterCollect map[string]*Router

// 添加路由容器
func (this RouterCollect) Add(relativePath string, router *Router) *RouterCollect {
	this[relativePath] = router
	return &this
}

// 所有路由
func (this RouterCollect) All() RouterCollect {
	return this
}

// 获取路由
func (this RouterCollect) Get(relativePath string) *Router {
	router, ok := this[relativePath]
	if !ok {
		throw.Throw(relativePath+"未定义路由", 0)
	}
	return router
}

func (this *Router) Method() string {
	return this.method
}

func (this *Router) RelativePath() string {
	return this.relativePath
}

func (this *Router) Flag() string {
	return this.flag
}

func (this *Router) FrontPath() string {
	return this.frontPath
}

func (this *Router) Handle() interface{} {
	return this.handle
}

func (this *Router) IsStatic() bool {
	return this.isStatic
}

func (this *Router) Desc() string {
	return this.desc
}

func (this *Router) attribute(args ...*RouteAttribute)  {
	for _, attribute := range args {
		if attribute.Name == ROUTE_FLAG {
			this.flag = attribute.Value.(string)
		} else if attribute.Name == ROUTE_FRONTPATH {
			this.frontPath = attribute.Value.(string)
		} else if attribute.Name == ROUTE_DESC {
			this.desc = attribute.Value.(string)
		} else if attribute.Name == ROUTE_IS_STATIC {
			this.isStatic = attribute.Value.(bool)
		}
	}
}

func (this *Router) Get(relativePath string, handle interface{}, args ...*RouteAttribute) *Router {
	this.method = "GET"
	this.relativePath = relativePath
	this.handle = handle
	this.attribute(args...)
	return this
}

func (this *Router) Post(relativePath string, handle interface{}, args ...*RouteAttribute) *Router {
	this.method = "POST"
	this.relativePath = relativePath
	this.handle = handle
	this.attribute(args...)
	return this
}

func (this *Router) Put(relativePath string, handle interface{}, args ...*RouteAttribute) *Router {
	this.method = "PUT"
	this.relativePath = relativePath
	this.handle = handle
	this.attribute(args...)
	return this
}

func (this *Router) Delete(relativePath string, handle interface{}, args ...*RouteAttribute) *Router {
	this.method = "DELETE"
	this.relativePath = relativePath
	this.handle = handle
	this.attribute(args...)
	return this
}
