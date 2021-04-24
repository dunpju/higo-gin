package higo

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Websocket interface{}

type WebsocketPingFunc func(websocketConn *WebsocketConn, waittime time.Duration)

type WebsocketClient struct {
	clients sync.Map
}

func NewWebsocketClient() *WebsocketClient {
	return &WebsocketClient{}
}

func (this *WebsocketClient) Store(conn *websocket.Conn) {
	wsConn := NewWebsocketConn(conn)
	this.clients.Store(conn.RemoteAddr().String(), wsConn)
	go wsConn.Ping(time.Second * 1) //心跳
	go wsConn.WriteLoop()//写循环
	go wsConn.ReadLoop()//读循环
	go wsConn.HandlerLoop()//处理控制循环
}

func (this *WebsocketClient) SendAll(msg string) {
	this.clients.Range(func(key, client interface{}) bool {
		conn := client.(*WebsocketConn).conn
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			//TODO::应该记录日志
			this.Remove(conn)
			panic(err)
		}
		return true
	})
}

func (this *WebsocketClient) Remove(conn *websocket.Conn) {
	this.clients.Delete(conn.RemoteAddr().String())
}

//webSocket请求连接中间件
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
