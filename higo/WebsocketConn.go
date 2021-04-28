package higo

import (
	"bytes"
	"fmt"
	"github.com/dengpju/higo-router/router"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"time"
)

type WebsocketConn struct {
	route     *router.Route
	conn      *websocket.Conn
	readChan  chan *WebsocketMessage
	writeChan chan []byte
	closeChan chan byte
}

func NewWebsocketConn(route *router.Route, conn *websocket.Conn) *WebsocketConn {
	return &WebsocketConn{route: route, conn: conn, readChan: make(chan *WebsocketMessage),
		writeChan: make(chan []byte), closeChan: make(chan byte)}
}

func (this *WebsocketConn) Conn() *websocket.Conn {
	return this.conn
}

func (this *WebsocketConn) Ping(waittime time.Duration) {
	for {
		WebsocketPingHandler(this, waittime)
	}
}

func (this *WebsocketConn) ReadLoop() {
	for {
		t, message, err := this.conn.ReadMessage()
		if err != nil {
			this.conn.Close()
			WebsocketContainer.Remove(this.conn)
			this.closeChan <- 1
			break
		}
		this.readChan <- NewWebsocketMessage(t, message)
	}
}

func (this *WebsocketConn) WriteLoop() {
loop:
	for {
		select {
		case msg := <-this.writeChan:
			if err := this.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				this.conn.Close()
				WebsocketContainer.Remove(this.conn)
				this.closeChan <- 1
				break loop
			}
		}
	}
}

func (this *WebsocketConn) HandlerLoop() {
loop:
	for {
		select {
		case msg := <-this.readChan:
			//TODO::做路由转发
			fmt.Println("HandlerLoop", this.route)
			fmt.Println("HandlerLoop", string(msg.MessageData))
			fmt.Println("HandlerLoop", this.conn.RemoteAddr().String())

			handle := this.route.Handle()
			ctx := &gin.Context{Request: &http.Request{PostForm: make(url.Values)}}
			reader := bytes.NewReader(msg.MessageData)
			request,_ := http.NewRequest(router.POST, this.route.FullPath(), reader)
			request.Header.Set("Content-Type", "application/json")
			ctx.Request = request

			//调度
			responser := handle.(func(*gin.Context) Websocket)(ctx)
			// 写数据
			this.writeChan <- []byte("receiv: " + responser.(string))
		case <-this.closeChan:
			fmt.Println("已经关闭")
			break loop
		}
	}
}

func GetWebsocketConn(ctx *gin.Context) *WebsocketConn {
	ip, ok := ctx.Get(WsConnIp)
	if !ok {
		panic("websocket conn ip non-existent")
	}
	if conn, ok := WebsocketContainer.clients.Load(ip); ok {
		return conn.(*WebsocketConn)
	} else {
		panic("websocket conn non-existent")
	}
}

//webSocket请求连接
func websocketConnFunc(ctx *gin.Context) string {
	//升级get请求为webSocket协议
	client, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("websocketConnFunc Method", ctx.Request.Method)
	fmt.Println("websocketConnFunc URL", ctx.Request.URL)
	route := router.GetRoutes(WebsocketServe).Route(ctx.Request.Method, ctx.Request.URL.Path).SetHeader(ctx.Request.Header)
	fmt.Println("websocketConnFunc route", route)
	WebsocketContainer.Store(route, client)
	return client.RemoteAddr().String()
}
