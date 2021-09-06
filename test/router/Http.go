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
type Http struct {
	*higo.Serve `inject:"Bean.NewServe('env.serve.HTTP_HOST')"`
}

func NewHttp() *Http {
	return &Http{}
}

func (this *Http) Loader(hg *higo.Higo) {

	// 静态文件
	hg.StaticFile("/", fmt.Sprintf("%sdist", hg.GetRoot().Join(utils.PathSeparator())))
	this.http(hg)
}

// api 路由
func (this *Http) http(hg *higo.Higo) {
	hg.Get("/http/test_throw", Controllers.HttpsTestThrow, hg.Flag("TestThrow"), hg.Desc("测试异常"))
	hg.Get("/http/test_get", Controllers.HttpsTestGet, hg.Flag("TestGet"), hg.Desc("测试GET"))
	hg.Post("/http/test_post", Controllers.HttpsTestPost, hg.Flag("TestPost"), hg.Desc("测试POST"))
	// 路由组
	hg.AddGroup("/http/v2", func() {
		hg.Get("/test_throw", V2.HttpsTestThrow, hg.Flag("TestThrow"), hg.Desc("v2 测试异常"))
		hg.Get("/test_get", V2.HttpsTestGet, hg.Flag("TestGet"), hg.Desc("v2 测试GET"))
		hg.Post("/test_post", V2.HttpsTestPost, hg.Flag("TestPost"), hg.Desc("v2 测试POST"))
	})
	hg.AddGroup("/http/v3", func() {
		hg.AddGroup("/user", func() {
			hg.Post("/login", V3.NewDemoController().Login, hg.Flag("Login"), hg.Desc("V3 登录"))
		})
		hg.Get("/test_throw", V3.NewDemoController().HttpsTestThrow, hg.Flag("TestThrow"), hg.Desc("V3 测试异常"))
		hg.Get("/test_get", V3.NewDemoController().HttpsTestGet, hg.Flag("TestGet"), hg.Desc("V3 测试GET111"))
		hg.Post("/test_post", V3.NewDemoController().HttpsTestPost, hg.Flag("TestPost"), hg.Desc("V3 测试POST"))
		hg.Get("/test_get_redis", V3.NewRedisController().Test, hg.Flag("test_get_redis"), hg.Desc("V3 测试redis"))

	})
}
