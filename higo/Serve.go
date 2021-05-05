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

func newServe() *Serve {
	return &Serve{Middle: make([]IMiddleware, 0)}
}

func NewServe(conf string) *Serve {
	serve := newServe()
	if "" != conf {
		configs := config.Get(conf).(*config.Configure)
		name := configs.Get("Name").(string)
		t := configs.Get("Type").(string)
		serve.Name = name
		serve.Type = t
		serve.Config = conf
	}
	return serve
}

func (this *Serve) SetRouter(router IRouterLoader) *Serve {
	this.Router = router
	return this
}

func (this *Serve) SetMiddle(middles ...IMiddleware) *Serve {
	this.Middle = append(this.Middle, middles...)
	return this
}

func (this *Serve) GetServe() *Serve {
	return this
}

func (this *Serve) Loader(hg *Higo) *Higo {

	return hg
}
