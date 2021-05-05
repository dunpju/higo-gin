package router

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers"
)

// https api 接口
type Websocket struct {
	*higo.Serve `inject:"Bean.NewServe('env.serve.WEBSOCKET_HOST')"`
}

func NewWebsocket() *Websocket {
	return &Websocket{}
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
