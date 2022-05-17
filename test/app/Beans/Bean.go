package Beans

import (
	"github.com/dengpju/higo-gin/higo"
	a "github.com/dengpju/higo-gin/test/app/Controllers"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-gin/test/app/Services"
	"github.com/gomodule/redigo/redis"
)

type MyBean struct{ higo.Bean }

func NewMyBean() *MyBean {
	return &MyBean{}
}

func (this *MyBean) DemoService() *Services.DemoService {
	return Services.NewDemoService()
}

func (this *MyBean) NewOrm() *higo.Orm {
	return higo.NewOrm()
}

func (this *MyBean) NewRedisPool() *redis.Pool {
	return higo.RedisPool
}

func (this *MyBean) NewRedisAdapter() *higo.RedisAdapter {
	return higo.NewRedisAdapter()
}

func (this *MyBean) NewRedisController() *V3.RedisController {
	return V3.NewRedisController()
}

func (this *MyBean) NewDemoController() *V3.DemoController {
	return V3.NewDemoController()
}

func (this *MyBean) NewEventController() *a.EventController {
	return a.NewEventController()
}

func (this *MyBean) New_gen_github8com_dengpju_higo9gin_test_app_Controllers_YaaController() *a.YaaController {
	return a.NewYaaController()
}