package higo

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

type WebsocketConn struct {
	conn *websocket.Conn
}

func NewWebsocketConn(conn *websocket.Conn) *WebsocketConn {
	return &WebsocketConn{conn: conn}
}

func (this *WebsocketConn) Conn() *websocket.Conn {
	return this.conn
}

func (this *WebsocketConn) Ping(waittime time.Duration) {
	for {
		WebsocketPingHandler(this, waittime)
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
