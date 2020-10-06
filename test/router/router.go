package router

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

// https api 接口
type Https struct {}

func NewHttps() *Https  {
	return &Https{}
}

// 路由装载器
func (this *Https) Loader(hg *higo.Higo) *higo.Higo {

	// 静态文件
	hg.Engine.StaticFile("/", fmt.Sprintf("%sdist", hg.GetRoot()))
	this.Api(hg)

	return hg
}

// api 路由(需要首先在容器中添加Map映射)
func (this *Https) Api(hg *higo.Higo) {
	hg.AddRoute(
		higo.Route{Method: "GET", RelativePath: "/test_get", Handle: test_get, Flag: "test_get"},
		higo.Route{Method: "GET", RelativePath: "/test_throw", Handle: test_throw, Flag: "test_throw"},
		higo.Route{Method: "post", RelativePath: "/test_post", Handle: test_post, Flag: "test_post"},
	)
}

// 测试异常
func test_throw(ctx *gin.Context) string  {
	higo.Throw("测试异常", 0)
	return "test_throw"
}

// 测试get请求
func test_get(ctx *gin.Context) string  {
	return "test_get"
}

// 测试post请求
func test_post(ctx *gin.Context) string {
	return "test_post"
}