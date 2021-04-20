package higo

import (
	"github.com/dengpju/higo-config/config"
)

type Serve struct {
	Name   string
	Type   string
	Config string
	Router IRouterLoader
}

func NewServe(conf string, router IRouterLoader) *Serve {
	configs := config.Get(conf).(config.Configure)
	name := configs.Get("Name").(string)
	t := configs.Get("Type").(string)
	return &Serve{Name: name, Type: t, Config: conf, Router: router}
}
