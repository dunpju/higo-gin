package higo

import (
	"github.com/dunpju/higo-router/router"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type WebsocketCheckFunc func(r *http.Request) bool

type WebsocketPingFunc func(websocketConn *WebsocketConn, wait time.Duration)

type WebsocketClient struct {
	clients sync.Map
}

func NewWebsocketClient() *WebsocketClient {
	return &WebsocketClient{}
}

func (this *WebsocketClient) Store(route *router.Route, conn *websocket.Conn) {
	wsConn := NewWebsocketConn(route, conn)
	this.clients.Store(conn.RemoteAddr().String(), wsConn)
	go wsConn.Ping(WsPitpatSleep) //心跳
	go wsConn.WriteLoop()         //写循环
	go wsConn.ReadLoop()          //读循环
	go wsConn.HandlerLoop()       //处理控制循环
}

func (this *WebsocketClient) SendAll(msg string) {
	this.clients.Range(func(key, client interface{}) bool {
		conn := client.(*WebsocketConn).conn
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			this.Remove(conn)
		}
		return true
	})
}

func (this *WebsocketClient) Remove(conn *websocket.Conn) {
	this.clients.Delete(conn.RemoteAddr().String())
}

// ws连接中间件
func wsConnMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn := websocketConnFunc(ctx)
		// 设置变量到Context的key中，可以通过Get()取
		ctx.Set(WsConnIp, conn)
		// 执行函数
		ctx.Next()
		// 中间件执行完后续的一些事情
		//status := ctx.Writer.Status()
	}
}

// 连接升级协议handle
func wsUpGraderHandle(route *router.Route) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := ctx.Get(WsConnIp)
		if !ok {
			panic("websocket conn ip non-existent")
		}
	}
}
