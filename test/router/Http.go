package router

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers"
	"github.com/dengpju/higo-gin/test/app/Controllers/V2"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
)

// https api 接口
type Http struct {}

func NewHttp() *Http  {
	return &Http{}
}

func (this *Http) Loader(hg *higo.Higo) *higo.Higo {

	// 静态文件
	hg.StaticFile("/", fmt.Sprintf("%sdist", hg.GetRoot()))
	this.http(hg)

	return hg
}

// api 路由
func (this *Http) http(hg *higo.Higo) {
	hg.AddRoute(
		higo.Route{Method: "GET", RelativePath: "/test_throw", Handle: Controllers.HttpTestThrow, Flag: "TestThrow", Desc:"测试异常"},
		higo.Route{Method: "GET", RelativePath: "/test_get", Handle: Controllers.HttpTestGet, Flag: "TestGet", Desc:"测试GET"},
		higo.Route{Method: "post", RelativePath: "/test_post", Handle: Controllers.HttpTestPost, Flag: "TestPost", Desc:"测试POST"},
	)
	// 路由组
	hg.AddGroup("v2",
		higo.Route{Method: "GET", RelativePath: "/test_throw", Handle: V2.HttpTestThrow, Flag: "TestThrow", Desc:"V2 测试异常"},
		higo.Route{Method: "GET", RelativePath: "/test_get", Handle: V2.HttpTestGet, Flag: "TestGet", Desc:"V2 测试GET"},
		higo.Route{Method: "post", RelativePath: "/test_post", Handle: V2.HttpTestPost, Flag: "TestPost", Desc:"V2 测试POST"},
	)
	// 路由组
	hg.AddGroup("v3",
		higo.Route{Method: "GET", RelativePath: "/test_throw", Handle: V3.NewDemoController().HttpsTestThrow, Flag: "TestThrow", Desc:"V3 测试异常"},
		higo.Route{Method: "GET", RelativePath: "/test_get", Handle: V3.NewDemoController().HttpsTestGet, Flag: "TestGet", Desc:"V3 测试GET"},
		higo.Route{Method: "post", RelativePath: "/test_post", Handle: V3.NewDemoController().HttpsTestPost, Flag: "TestPost", Desc:"V3 测试POST"},
	)
}