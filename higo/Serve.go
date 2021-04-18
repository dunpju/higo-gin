package higo

import "gitee.com/dengpju/higo-configure/configure"

type Serve struct {
	Name   string
	Type   string
	Config string
	Router IRouterLoader
}

func NewServe(config string, router IRouterLoader) *Serve {
	configs := configure.Config(config)
	name := configs.String("Name")
	t := configs.String("Type")
	return &Serve{Name: name, Type: t, Config: config, Router: router}
}
