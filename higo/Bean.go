package higo

import (
	"github.com/gomodule/redigo/redis"
)

type Bean struct {
	Middleware
}

func NewBean() *Bean {
	return &Bean{}
}

func (this *Bean) Provider() {}

func (this *Bean) NewServe(conf string) *Serve {
	return NewServe(conf)
}

func (this *Bean) NewOrm() *Orm {
	return NewOrm()
}

func (this *Bean) NewRedisPool() *redis.Pool {
	return RedisPool
}

func (this *Bean) NewRedisAdapter() *RedisAdapter {
	return NewRedisAdapter()
}
