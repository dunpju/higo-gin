package Config

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Services"
	"github.com/gomodule/redigo/redis"
)

type MyBean struct {
	higo.Bean
}

func NewMyBean() *MyBean {
	return &MyBean{}
}

func (this *MyBean) DemoService() *Services.DemoService {
	return Services.NewDemoService()
}

func (this *MyBean) NewGorm() *higo.Gorm {
	return higo.NewGorm()
}

func (this *MyBean) NewRedisPool() *redis.Pool {
	return higo.RedisPool
}

func (this *MyBean) NewRedisAdapter() *higo.RedisAdapter {
	return higo.NewRedisAdapter()
}
