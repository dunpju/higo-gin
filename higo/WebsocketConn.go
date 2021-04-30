package higo

import (
	"bytes"
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-logger/logger"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/exception"
	"github.com/dengpju/higo-utils/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	//Recover处理函数(可自定义替换)
	WsRecoverHandle WsRecoverFunc
	wsRecoverOnce   sync.Once
)

func init() {
	wsRecoverOnce.Do(func() {
		WsRecoverHandle = func(r interface{}) (respMsg string) {
			if msg, ok := r.(*code.Code); ok {
				respMsg = make(utils.MapString).
					Put("code", msg.Code).
					Put("message", msg.Message).
					Put("data", nil).
					String()
			} else if MapString, ok := r.(utils.MapString); ok {
				respMsg = MapString.String()
			} else {
				respMsg = make(utils.MapString).
					Put("code", 0).
					Put("message", exception.ErrorToString(r)).
					Put("data", nil).
					String()
			}
			return
		}
	})
}

type WsRecoverFunc func(r interface{}) string

type WebsocketConn struct {
	route     *router.Route
	conn      *websocket.Conn
	readChan  chan *WsReadMessage
	writeChan chan WsWriteMessage
	closeChan chan byte
}

func NewWebsocketConn(route *router.Route, conn *websocket.Conn) *WebsocketConn {
	return &WebsocketConn{route: route, conn: conn, readChan: make(chan *WsReadMessage),
		writeChan: make(chan WsWriteMessage), closeChan: make(chan byte)}
}

func (this *WebsocketConn) Conn() *websocket.Conn {
	return this.conn
}

func (this *WebsocketConn) Ping(waittime time.Duration) {
	for {
		WsPingHandle(this, waittime)
	}
}

func (this *WebsocketConn) ReadLoop() {
	for {
		t, message, err := this.conn.ReadMessage()
		if err != nil {
			this.Close()
			break
		}
		this.readChan <- NewReadMessage(t, message)
	}
}

func (this *WebsocketConn) Close() {
	this.conn.Close()
	WsContainer.Remove(this.conn)
	this.closeChan <- 1
}

func (this *WebsocketConn) WriteLoop() {
loop:
	for {
		select {
		case msg := <-this.writeChan:
			if WsResperror == msg.MessageType {
				_ = this.conn.WriteMessage(websocket.TextMessage, msg.MessageData)
				this.Close()
				break loop
			}
			if err := this.conn.WriteMessage(websocket.TextMessage, msg.MessageData); err != nil {
				this.Close()
				break loop
			}
		}
	}
}

func (this *WebsocketConn) HandlerLoop() {
	defer func() {
		if r := recover(); r != nil {
			logger.LoggerStack(r, utils.GoroutineID())
			this.writeChan <- WsRespError(WsRecoverHandle(r))
		}
	}()
loop:
	for {
		select {
		case msg := <-this.readChan:
			// 写数据
			this.writeChan <- this.dispatch(msg)
		case <-this.closeChan:
			logger.Logrus.Info("websocket conn " + this.Conn().RemoteAddr().String() + " have already closed")
			break loop
		}
	}
}

func (this *WebsocketConn) dispatch(msg *WsReadMessage) WsWriteMessage {
	handle := this.route.Handle()
	ctx := &gin.Context{Request: &http.Request{PostForm: make(url.Values)}}
	reader := bytes.NewReader(msg.MessageData)
	request, _ := http.NewRequest(router.POST, this.route.FullPath(), reader)
	request.Header.Set("Content-Type", "application/json")
	ctx.Request = request

	return handle.(func(*gin.Context) WsWriteMessage)(ctx)
}

func WsConn(ctx *gin.Context) *WebsocketConn {
	client, ok := ctx.Get(WsConnIp)
	if !ok {
		panic("websocket conn client non-existent")
	}
	if conn, ok := WsContainer.clients.Load(client); ok {
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

	route := router.GetRoutes(WebsocketServe).Route(ctx.Request.Method, ctx.Request.URL.Path).SetHeader(ctx.Request.Header)

	WsContainer.Store(route, client)
	return client.RemoteAddr().String()
}
