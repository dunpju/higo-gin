package V3

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/higo/injector"
	"github.com/dengpju/higo-gin/test/app/Exception"
	"github.com/gin-gonic/gin"
)

type DemoController struct {
	*higo.HgController
	Age *injector.Value `prefix:"user.age"`
}

func NewDemoController(hg *higo.Higo) *DemoController {
	return &DemoController{}
}

// 测试异常
func (this *DemoController) HttpsTestThrow(ctx *gin.Context) string {
	fmt.Println(ctx.Query("id"))
	fmt.Println(111)
	fmt.Println(this.Higo.ApplicationContext())
	fmt.Println(this.Age.String())
	var s []map[string]interface{}
	m1 := make(map[string]interface{})
	m1["jj"] = "m1jjj"
	m1["dd"] = "m1ddd"
	m2 := make(map[string]interface{})
	m2["jj"] = "m2jjj"
	m2["dd"] = "m2ddd"
	s = append(s, m1)
	s = append(s, m2)
	Exception.NewBusinessException(2, "v3 https 测试异常", s)
	higo.Throw("v3 https 测试异常", 2, struct {
		Id   int
		Name string
	}{Id: 1, Name: "哦哦"})
	return "v3 https_test_throw"
}

// 测试get请求
func (this *DemoController) HttpsTestGet(ctx *gin.Context) string  {
	fmt.Printf("%p\n", this)
	return "v3 https_test_get"
}

// 测试post请求
func (this *DemoController) HttpsTestPost(ctx *gin.Context) string {
	return "v3 https_test_post"
}

// 测试异常
func (this *DemoController) HttpTestThrow(ctx *gin.Context) string  {
	higo.Throw("v3 http 测试异常", 0)
	return "v3 http_test_throw"
}

// 测试get请求
func (this *DemoController) HttpTestGet(ctx *gin.Context) string  {
	return "v3 http_test_get"
}

// 测试post请求
func (this *DemoController) HttpTestPost(ctx *gin.Context) string {
	return "v3 http_test_post"
}
