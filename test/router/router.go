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
	hg.Engine.StaticFile("/", fmt.Sprintf(".%sdist", higo.PathSeparator))
	this.Api(hg)

	return hg
}

// api 路由(需要首先在容器中添加Map映射)
func (this *Https) Api(hg *higo.Higo) {
	hg.AddRoute(
		higo.Route{Method: "GET", RelativePath: "/test_get", Handle: test_get, Flag: "AttackList"},
		higo.Route{Method: "post", RelativePath: "/test_post", Handle: test_post, Flag: "Login"},
	)
}

// 测试get请求
func test_get(ctx *gin.Context) string  {
	return "test_get"
}

// 测试post请求
func test_post(ctx *gin.Context) string {
	return "test_post"
}