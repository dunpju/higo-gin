package higo

import (
	"github.com/dengpju/higo-config/config"
)

type Serve struct {
	Name   string
	Type   string
	Config string
	Router IRouterLoader
	Middle []IMiddleware
}

func NewServe(conf string, router IRouterLoader, middles ...IMiddleware) *Serve {
	configs := config.Get(conf).(*config.Configure)
	name := configs.Get("Name").(string)
	t := configs.Get("Type").(string)
	serve := &Serve{Name: name, Type: t, Config: conf, Router: router, Middle: make([]IMiddleware, 0)}
	if len(middles) > 0 {
		serve.Middle = middles
	}
	return serve
}
