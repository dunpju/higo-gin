package router

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers"
	"github.com/dengpju/higo-gin/test/app/Controllers/V2"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-utils/utils"
)

// https api 接口
type Https struct {
}

func NewHttps() *Https {
	return &Https{}
}

func (this *Https) Serve(middles ...higo.IMiddleware) *higo.Serve {
	return higo.NewServe("env.serve.HTTPS_HOST", this, middles...)
}

// 路由装载器
func (this *Https) Loader(hg *higo.Higo) *higo.Higo {

	// 静态文件
	hg.StaticFile("/", fmt.Sprintf("%sdist", hg.GetRoot().Separator(utils.PathSeparator())))
	this.Api(hg)

	return hg
}

// api 路由
func (this *Https) Api(hg *higo.Higo) {
	hg.Get("/https/test_throw", Controllers.HttpsTestThrow, hg.Flag("TestThrow"), hg.Desc("测试异常"))
	hg.Get("/https/test_get", Controllers.HttpsTestGet, hg.Flag("TestGet"), hg.Desc("测试GET"))
	hg.Post("/https/test_post", Controllers.HttpsTestPost, hg.Flag("TestPost"), hg.Desc("测试POST"))
	// 路由组
	hg.AddGroup("/https/v2", func() {
		hg.Get("/test_throw", V2.HttpsTestThrow, hg.Flag("TestThrow"), hg.Desc("v2 测试异常"))
		hg.Get("/test_get", V2.HttpsTestGet, hg.Flag("TestGet"), hg.Desc("v2 测试GET"))
		hg.Post("/test_post", V2.HttpsTestPost, hg.Flag("TestPost"), hg.Desc("v2 测试POST"))
	})
	hg.AddGroup("/https/v3", func() {
		hg.AddGroup("/user", func() {
			hg.Post("/login", V3.NewDemoController().Login, hg.Flag("Login"), hg.Desc("V3 登录"))
		})
		hg.Get("/test_throw", V3.NewDemoController().HttpsTestThrow, hg.Flag("TestThrow"), hg.Desc("V3 测试异常"))
		hg.Get("/test_get", V3.NewDemoController().HttpsTestGet, hg.Flag("TestGet"), hg.Desc("V3 测试GET"))
		hg.Post("/test_post", V3.NewDemoController().HttpsTestPost, hg.Flag("TestPost"), hg.Desc("V3 测试POST"))
		hg.Get("/test_get_redis", V3.NewRedisController().Test, hg.Flag("test_get_redis"), hg.Desc("V3 测试redis"))

	})

	hg.Route(V3.NewRedisController())
}
