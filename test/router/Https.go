package router

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers"
	"github.com/dengpju/higo-gin/test/app/Controllers/V2"
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

// api 路由
func (this *Https) Api(hg *higo.Higo) {
	hg.AddRoute(
		higo.Route{Method: "GET", RelativePath: "/test_throw", Handle: Controllers.HttpsTestThrow, Flag: "TestThrow", Desc:"测试异常"},
		higo.Route{Method: "GET", RelativePath: "/test_get", Handle: Controllers.HttpsTestGet, Flag: "TestGet", Desc:"测试GET"},
		higo.Route{Method: "post", RelativePath: "/test_post", Handle: Controllers.HttpsTestPost, Flag: "TestPost", Desc:"测试POST"},
	)
	// 路由组
	hg.AddGroup("v2",
		higo.Route{Method: "GET", RelativePath: "/test_throw", Handle: V2.HttpsTestThrow, Flag: "TestThrow", Desc:"V2 测试异常"},
		higo.Route{Method: "GET", RelativePath: "/test_get", Handle: V2.HttpsTestGet, Flag: "TestGet", Desc:"V2 测试GET"},
		higo.Route{Method: "post", RelativePath: "/test_post", Handle: V2.HttpsTestPost, Flag: "TestPost", Desc:"V2 测试POST"},
	)
}