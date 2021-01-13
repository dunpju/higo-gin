package router

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers"
	"github.com/dengpju/higo-gin/test/app/Controllers/V2"
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
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_throw"), higo.Handle(Controllers.HttpTestThrow), higo.Flag("TestThrow"), higo.Desc("测试异常")),
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_get"), higo.Handle(Controllers.HttpTestGet), higo.Flag("TestGet"), higo.Desc("测试GET")),
		higo.Route(higo.Method("post"), higo.RelativePath("/test_post"), higo.Handle(Controllers.HttpTestPost), higo.Flag("TestPost"), higo.Desc("测试POST")),
	)
	// 路由组
	hg.AddGroup("v2",
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_throw"), higo.Handle(V2.HttpTestThrow), higo.Flag("TestThrow"), higo.Desc("V2 测试异常")),
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_get"), higo.Handle(V2.HttpTestGet), higo.Flag("TestGet"), higo.Desc("V2 测试GET")),
		higo.Route(higo.Method("post"), higo.RelativePath("/test_post"), higo.Handle(V2.HttpTestPost), higo.Flag("TestPost"), higo.Desc("V2 测试POST")),
	)
	// 路由组
	hg.AddGroup("v3",
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_throw"), higo.Handle(V2.HttpTestThrow), higo.Flag("TestThrow"), higo.Desc("V3 测试异常")),
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_get"), higo.Handle(V2.HttpTestGet), higo.Flag("TestGet"), higo.Desc("V3 测试GET")),
		higo.Route(higo.Method("post"), higo.RelativePath("/test_post"), higo.Handle(V2.HttpTestPost), higo.Flag("TestPost"), higo.Desc("V3 测试POST")),
	)
}