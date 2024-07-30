package Middlewares

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-wsock/wsock"
	"github.com/gin-gonic/gin"
)

// Websocket 运行日志
type Websocket struct{}

// NewWebsocket 构造函数
func NewWebsocket() *Websocket {
	return &Websocket{}
}

func (this *Websocket) Middle(hg *higo.Higo) gin.HandlerFunc {
	wsock.SetServe("http")
	return wsock.ConnUpGrader()
}
