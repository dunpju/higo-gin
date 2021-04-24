package higo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

type WebsocketConn struct {
	conn      *websocket.Conn
	readChan  chan *WebsocketMessage
	writeChan chan []byte
	closeChan chan byte
}

func NewWebsocketConn(conn *websocket.Conn) *WebsocketConn {
	return &WebsocketConn{conn: conn, readChan: make(chan *WebsocketMessage), writeChan: make(chan []byte), closeChan: make(chan byte)}
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
			fmt.Println(string(msg.MessageData))
			// 写数据
			this.writeChan <- []byte("receiv: " + string(msg.MessageData))
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
	WebsocketContainer.Store(client)
	return client.RemoteAddr().String()
}
