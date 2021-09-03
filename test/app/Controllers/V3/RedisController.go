package V3

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

type RedisController struct {
	*higo.Orm `inject:"Bean.NewOrm()"`
	Redis      *higo.RedisAdapter `inject:"Bean.NewRedisAdapter()"`
}

func NewRedisController() *RedisController {
	return &RedisController{}
}

func (this *RedisController) New() higo.IClass {
	return NewRedisController()
}

func (this *RedisController) Route(hg *higo.Higo) {
	hg.AddGroup("/https/v3", func() {
		hg.Get("/test_get_redis", this.Test, hg.Flag("test_get_redis"), hg.Desc("V3 测试redis"))
	})

	// 路由组
	hg.AddGroup("/https/v4", func() {
		hg.Get("/test_redis", this.Test, hg.Flag("TestThrow"), hg.Desc("V4 测试redis"))
		hg.AddGroup("/v5", func() {
			hg.Get("/get_test_redis", this.Test, hg.Flag("get_test_redis"), hg.Middle(this.MiddleWare()))
		}, hg.GroupMiddle(this.V5GroupMiddleWare()))
	}, hg.GroupMiddle(this.V4GroupMiddleWare()))
}

func (this *RedisController) Test(ctx *gin.Context) string {
	ctx.Set("db_result", rand.Intn(1000))
	higo.Redis.Set("name", rand.Intn(1000))
	v := higo.Redis.Get("name")
	return v
}

func (this *RedisController) V4GroupMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Println("V4 Group中间件开始执行了")
		// 设置变量到Context的key中，可以通过Get()取
		c.Set("request", "V4 Group中间件")
		// 执行函数
		c.Next()
		// 中间件执行完后续的一些事情
		status := c.Writer.Status()
		fmt.Println("V4 Group中间件执行完毕", status)
		t2 := time.Since(t)
		fmt.Println("V4 Group time:", t2)
	}
}

func (this *RedisController) V5GroupMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Println("V5 Group中间件开始执行了")
		// 设置变量到Context的key中，可以通过Get()取
		c.Set("request", "V5 Group中间件")
		// 执行函数
		c.Next()
		// 中间件执行完后续的一些事情
		status := c.Writer.Status()
		fmt.Println("V5 Group中间件执行完毕", status)
		t2 := time.Since(t)
		fmt.Println("V5 Group time:", t2)
	}
}

func (this *RedisController) MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Println("中间件开始执行了")
		// 设置变量到Context的key中，可以通过Get()取
		c.Set("request", "中间件")
		// 执行函数
		c.Next()
		// 中间件执行完后续的一些事情
		status := c.Writer.Status()
		fmt.Println("中间件执行完毕", status)
		t2 := time.Since(t)
		fmt.Println("time:", t2)
	}
}
