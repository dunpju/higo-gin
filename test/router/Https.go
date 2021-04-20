package router

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers"
	"github.com/dengpju/higo-gin/test/app/Controllers/V2"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-router/router"
)

// https api 接口
type Https struct {
}

func NewHttps() *Https {
	return &Https{}
}

func (this *Https) Serve() *higo.Serve {
	return higo.NewServe("app.HTTPS_HOST", this)
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
	router.Get("/test_throw", Controllers.HttpsTestThrow, router.Flag("TestThrow"), router.Desc("测试异常"))
	router.Get("/test_get", Controllers.HttpsTestGet, router.Flag("TestGet"), router.Desc("测试GET"))
	router.Post("/test_post", Controllers.HttpsTestPost, router.Flag("TestPost"), router.Desc("测试POST"))
	// 路由组
	router.AddGroup("/v2", func() {
		router.Get("/test_throw", V2.HttpsTestThrow, router.Flag("TestThrow"), router.Desc("v2 测试异常"))
		router.Get("/test_get", V2.HttpsTestGet, router.Flag("TestGet"), router.Desc("v2 测试GET"))
		router.Post("/test_post", V2.HttpsTestPost, router.Flag("TestPost"), router.Desc("v2 测试POST"))
	})
	router.AddGroup("/v3", func() {
		router.AddGroup("/user", func() {
			router.Post("/login", V3.NewDemoController().Login, router.Flag("Login"), router.Desc("V3 登录"))
		})
		router.Get("/test_throw", V3.NewDemoController().HttpsTestThrow, router.Flag("TestThrow"), router.Desc("V3 测试异常"))
		router.Get("/test_get", V3.NewDemoController().HttpsTestGet, router.Flag("TestGet"), router.Desc("V3 测试GET"))
		router.Post("/test_post", V3.NewDemoController().HttpsTestPost, router.Flag("TestPost"), router.Desc("V3 测试POST"))
		router.Get("/test_get_redis", V3.NewRedisController().Test, router.Flag("test_get_redis"), router.Desc("V3 测试redis"))

	})
	V3.NewRedisController().Route(hg)
}
