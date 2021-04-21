package higo

import (
	"github.com/gin-gonic/gin"
)

type WebsocketPong string

var WebsocketPongFunc = websocketPongFunc

//webSocket请求ping 返回pong
func websocketPongFunc(ctx *gin.Context) WebsocketPong {
	//升级get请求为webSocket协议
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return ""
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
	return ""
}
