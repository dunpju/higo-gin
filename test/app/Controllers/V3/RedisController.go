package V3

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/gin-gonic/gin"
	"math/rand"
	"sync"
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
	hg.AddGroup("v4",
		higo.Route(higo.Method("GET"), higo.RelativePath("/test_redis"), higo.Handle(this.Test), higo.Flag("TestThrow"), higo.Desc("V4 测试redis")),
	)
	return hg
}
