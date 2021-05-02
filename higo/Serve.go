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

func NewServe(conf string, router IRouterLoader) {
	configs := config.Get(conf).(*config.Configure)
	name := configs.Get("Name").(string)
	t := configs.Get("Type").(string)
	serve := &Serve{Name: name, Type: t, Config: conf, Router: router, Middle: make([]IMiddleware, 0)}
	router.SetServe(serve)
}

func (this *Serve) SetMiddle(middles ...IMiddleware) {
	this.Middle = append(this.Middle, middles...)
}
