package router

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers"
	"github.com/dengpju/higo-gin/test/app/Controllers/V2"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
)

// https api 接口
type Https struct{}

func NewHttps() *Https {
	return &Https{}
}

// 路由装载器
func (this *Https) Loader(hg *higo.Higo) *higo.Higo {

	// 静态文件
	hg.StaticFile("/", fmt.Sprintf("%sdist", hg.GetRoot()))
	this.Api(hg)

	return hg
}

// api 路由
func (this *Https) Api(hg *higo.Higo) {
	hg.AddRoute(
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_throw"), higo.Handle(Controllers.HttpsTestThrow), higo.Flag("TestThrow"), higo.Desc("测试异常")),
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_get"), higo.Handle(Controllers.HttpsTestGet), higo.Flag("TestGet"), higo.Desc("测试GET")),
		higo.Route(higo.Method("post"), higo.RelativePath("/test_post"), higo.Handle(Controllers.HttpsTestPost), higo.Flag("TestPost"), higo.Desc("测试POST")),
	)
	// 路由组
	hg.AddGroup("v2",
		hg.AddGroup("user",
			higo.Route(higo.Method("GET"), higo.RelativePath("/test_throw"), higo.Handle(V2.HttpsTestThrow), higo.Flag("TestThrow"), higo.Desc("V2 测试异常")),
			),
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_throw"), higo.Handle(V2.HttpsTestThrow), higo.Flag("TestThrow"), higo.Desc("V2 测试异常")),
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_get"), higo.Handle(V2.HttpsTestGet), higo.Flag("TestGet"), higo.Desc("V2 测试GET")),
		higo.Route(higo.Method("post"), higo.RelativePath("/test_post"), higo.Handle(V2.HttpsTestPost), higo.Flag("TestPost"), higo.Desc("V2 测试POST")),
	)
	// 路由组
	hg.AddGroup("v3",
		higo.Route(higo.Method("post"), higo.RelativePath("/user/login"), higo.Handle(V3.NewDemoController().Login), higo.Flag("Login"), higo.Desc("V3 登录")),
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_throw"), higo.Handle(V3.NewDemoController().HttpsTestThrow), higo.Flag("TestThrow"), higo.Desc("V3 测试异常")),
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_get"), higo.Handle(V3.NewDemoController().HttpsTestGet), higo.Flag("TestGet"), higo.Desc("V3 测试GET")),
		higo.Route(higo.Method("post"), higo.RelativePath("/test_post"), higo.Handle(V3.NewDemoController().HttpsTestPost), higo.Flag("TestPost"), higo.Desc("V3 测试POST")),
		higo.Route(higo.Method("get"), higo.RelativePath("/test_get_redis"), higo.Handle(V3.NewRedisController().Test), higo.Flag("test_get_redis"), higo.Desc("V3 测试redis")),
	)
	V3.NewRedisController().Route(hg)
}
