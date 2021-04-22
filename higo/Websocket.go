package higo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
	"time"
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

func (this *WebsocketClient) Store(conn *websocket.Conn) {
	wsConn := NewWebsocketConn(conn)
	this.clients.Store(conn.RemoteAddr().String(), wsConn)
	go wsConn.Ping(time.Second * 2)
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

//webSocket请求连接
func websocketConnFunc(ctx *gin.Context) WebsocketPong {
	//升级get请求为webSocket协议
	client, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}
	WebsocketContainer.Store(client)
	return "ok"
}

//webSocket请求ping 返回pong
func websocketPongFunc(ctx *gin.Context) WebsocketPong {
	//升级get请求为webSocket协议
	client, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}
	WebsocketContainer.Store(client)
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

type WebsocketConn struct {
	conn *websocket.Conn
}

func NewWebsocketConn(conn *websocket.Conn) *WebsocketConn {
	return &WebsocketConn{conn: conn}
}

func (this *WebsocketConn) Ping(waittime time.Duration) {
	for {
		time.Sleep(waittime)
		err := this.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			fmt.Println(WebsocketContainer)
			fmt.Println(err)
			WebsocketContainer.Remove(this.conn)
			return
		}
	}
}
