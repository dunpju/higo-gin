package router

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers"
)

// https api 接口
type Websocket struct {
	*higo.Serve
}

func NewWebsocket() *Websocket {
	this := &Websocket{}
	higo.NewServe("env.serve.WEBSOCKET_HOST", this)
	return this
}

func (this *Websocket) SetServe(serve *higo.Serve) {
	this.Serve = serve
}

func (this *Websocket) GetServe() *higo.Serve {
	return this.Serve
}

func (this *Websocket) Loader(hg *higo.Higo) *higo.Higo {
	hg.Route(Controllers.NewWebsocketController())

	this.api(hg)

	return hg
}

// api 路由
func (this *Websocket) api(hg *higo.Higo) {
	// 写路由
}
