package V3

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-router/router"
	"github.com/gin-gonic/gin"
	"math/rand"
	"sync"
	"time"
)

type RedisController struct {
	*higo.Gorm `inject:"Bean.NewGorm()"`
	Redis      *higo.RedisAdapter `inject:"Bean.NewRedisAdapter()"`
}

var redisControllerOnce sync.Once
var RedisControllerPointer *RedisController

func NewRedisController() *RedisController {
	redisControllerOnce.Do(func() {
		RedisControllerPointer = &RedisController{}
		injector.BeanFactory.Apply(RedisControllerPointer)
		injector.BeanFactory.Set(RedisControllerPointer)
		higo.AddContainer(RedisControllerPointer)
	})
	return RedisControllerPointer
}

func (this *RedisController) Self() higo.IClass {
	return this
}

func (this *RedisController) Test(ctx *gin.Context) string {
	ctx.Set("db_result", rand.Intn(1000))
	this.Redis.Set("name", rand.Intn(1000))
	v, _ := this.Redis.Get("name")
	return v
}

func (this *RedisController) Route(hg *higo.Higo) *higo.Higo {
	// 路由组
	router.AddGroup("/v4", func() {
		router.Get("/test_redis", this.Test, router.Flag("TestThrow"), router.Desc("V4 测试redis"))
		router.AddGroup("/v5", func() {
			router.Get("/get_test_redis", this.Test, router.Flag("get_test_redis"),router.Middleware(this.MiddleWare()))
		},router.GroupMiddle(this.V5GroupMiddleWare()))
	},router.GroupMiddle(this.V4GroupMiddleWare()))
	return hg
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