package higo

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
)

type WebsocketPong string

type WebsocketPongFunc func(ctx *gin.Context) WebsocketPong

type WebsocketClient struct {
	heartbeat sync.Map
	clients   sync.Map
}

func NewWebsocketClient() *WebsocketClient {
	return &WebsocketClient{}
}

func (this *WebsocketClient) Store(key string, conn *websocket.Conn) {
	this.clients.Store(key, conn)
}

func (this *WebsocketClient) SendAll(msg string) {
	this.clients.Range(func(key, client interface{}) bool {
		err := client.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			//TODO::应该记录日志
			panic(err)
		}
		return true
	})
}

//webSocket请求连接
func websocketConnFunc(ctx *gin.Context) WebsocketPong {
	//升级get请求为webSocket协议
	client, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}
	WebsocketClientContainer.Store(client.RemoteAddr().String(), client)
	return "ok"
}

//webSocket请求ping 返回pong
func websocketPongFunc(ctx *gin.Context) WebsocketPong {
	//升级get请求为webSocket协议
	client, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}
	WebsocketClientContainer.Store(client.RemoteAddr().String(), client)
	defer client.Close()
	for {
		//读取ws中的数据
		mt, message, err := client.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//写入ws数据
		err = client.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
	return "ok"
}
