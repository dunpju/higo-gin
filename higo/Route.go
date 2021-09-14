package higo

import (
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/exception"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

//路由容器
var RouterContainer RouterCollect

type RouterCollect map[string]*router.Route

// 添加路由容器
func (this RouterCollect) Add(unique string, route *router.Route) *RouterCollect {
	this[unique] = route
	return &this
}

// 所有路由
func (this RouterCollect) All() RouterCollect {
	return this
}

func Unique(method, absolutePath string) string {
	return method + "@" + absolutePath
}

// 路由唯一标记
func (this RouterCollect) Unique(method, absolutePath string) string {
	return Unique(method, absolutePath)
}

// 获取路由
func (this RouterCollect) Get(method, relativePath string) *router.Route {
	route, ok := this[Unique(method, relativePath)]
	if !ok {
		exception.Throw(exception.Message(relativePath+" Undefined route"), exception.Code(0))
	}
	return route
}

//添加路由
func (this *Higo) AddRoute(httpMethod string, relativePath string, handler interface{}, attributes ...*router.RouteAttribute) *Higo {
	funcForPC := strings.Split(runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(), "/")
	endFuncForPC := funcForPC[len(funcForPC)-1:]
	typeRegexp := regexp.MustCompile(`.*\.\(\*.*\)`)
	regResult := typeRegexp.FindStringSubmatch(endFuncForPC[0])
	if len(regResult) > 0 {
		regStruct := regResult[0]
		regStructName := strings.TrimRight(regStruct, ")")
		regStructName = "*" + strings.Join(strings.Split(regStructName, "(*"), "")
		typeRegexp = regexp.MustCompile(`\)\..*\-fm`)
		regMethod := typeRegexp.FindStringSubmatch(endFuncForPC[0])
		if len(regMethod) > 0 {
			method := regMethod[0]
			method = strings.TrimRight(strings.TrimLeft(method, ")."), "-fm")
			if nil != Di(regStructName) {
				router.AddRoute(httpMethod, relativePath, NewDispatch(Di(regStructName), method).Call(handler), attributes...)
			} else {
				router.AddRoute(httpMethod, relativePath, handler, attributes...)
			}
		} else {
			router.AddRoute(httpMethod, relativePath, handler, attributes...)
		}
	} else {
		router.AddRoute(httpMethod, relativePath, handler, attributes...)
	}
	return this
}

//添加路由组
func (this *Higo) AddGroup(prefix string, callable interface{}, attributes ...*router.RouteAttribute) *Higo {
	router.AddGroup(prefix, callable, attributes...)
	return this
}

//websocket路由
func (this *Higo) Ws(relativePath string, handler interface{}, attributes ...*router.RouteAttribute) *Higo {
	this.AddRoute(router.WEBSOCKET, relativePath, handler, attributes...)
	return this
}

func (this *Higo) Get(relativePath string, handler interface{}, attributes ...*router.RouteAttribute) *Higo {
	this.AddRoute(router.GET, relativePath, handler, attributes...)
	return this
}

func (this *Higo) Post(relativePath string, handler interface{}, attributes ...*router.RouteAttribute) *Higo {
	this.AddRoute(router.POST, relativePath, handler, attributes...)
	return this
}

func (this *Higo) Put(relativePath string, handler interface{}, attributes ...*router.RouteAttribute) *Higo {
	this.AddRoute(router.PUT, relativePath, handler, attributes...)
	return this
}

func (this *Higo) Delete(relativePath string, handler interface{}, attributes ...*router.RouteAttribute) *Higo {
	this.AddRoute(router.DELETE, relativePath, handler, attributes...)
	return this
}

func (this *Higo) Patch(relativePath string, handler interface{}, attributes ...*router.RouteAttribute) *Higo {
	this.AddRoute(router.PATCH, relativePath, handler, attributes...)
	return this
}

func (this *Higo) Head(relativePath string, handler interface{}, attributes ...*router.RouteAttribute) *Higo {
	this.AddRoute(router.HEAD, relativePath, handler, attributes...)
	return this
}

func (this *Higo) Flag(value string) *router.RouteAttribute {
	return router.NewRouteAttribute(router.RouteFlag, value)
}

func (this *Higo) FrontPath(value string) *router.RouteAttribute {
	return router.NewRouteAttribute(router.RouteFrontpath, value)
}

func (this *Higo) IsStatic(value bool) *router.RouteAttribute {
	return router.NewRouteAttribute(router.RouteIsStatic, value)
}

//描述
func (this *Higo) Desc(value string) *router.RouteAttribute {
	return router.NewRouteAttribute(router.RouteDesc, value)
}

//单路由中间件
func (this *Higo) Middle(value gin.HandlerFunc) *router.RouteAttribute {
	return router.NewRouteAttribute(router.RouteMiddleware, value)
}

//组中间件
func (this *Higo) GroupMiddle(value gin.HandlerFunc) *router.RouteAttribute {
	return router.NewRouteAttribute(router.RouteGroupMiddle, value)
}

func (this *Higo) SetServe(value interface{}) *router.RouteAttribute {
	return router.NewRouteAttribute(router.RouteServe, value)
}

func (this *Higo) SetHeader(value http.Header) *router.RouteAttribute {
	return router.NewRouteAttribute(router.RouteHeader, value)
}
