package Controllers

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-router/router"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WebsocketController struct {
	*higo.Gorm `inject:"Bean.NewGorm()"`
	Redis      *higo.RedisAdapter `inject:"Bean.NewRedisAdapter()"`
}

var (
	redisControllerOnce        sync.Once
	WebsocketControllerPointer *WebsocketController
)

func NewWebsocketController() *WebsocketController {
	redisControllerOnce.Do(func() {
		WebsocketControllerPointer = &WebsocketController{}
		injector.BeanFactory.Apply(WebsocketControllerPointer)
		injector.BeanFactory.Set(WebsocketControllerPointer)
		higo.AddContainer(WebsocketControllerPointer)
		higo.Upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	})
	return WebsocketControllerPointer
}

func (this *WebsocketController) Self() higo.IClass {
	return this
}

func (this *WebsocketController) Route(hg *higo.Higo) *higo.Higo {
	router.Get("/conn", this.Conn, router.Flag("WebsocketController.Conn"), router.Desc("conn"))
	router.Get("/send_all", this.SendAll, router.Flag("WebsocketController.SendAll"), router.Desc("SendAll"))
	return hg
}

//webSocket请求
func (this *WebsocketController) Conn(ctx *gin.Context) higo.Websocket {
fmt.Println("控制器")
	//ws := higo.GetWebsocketConn(ctx)
	/**
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
	 */

	/**
	//读取ws中的数据
	mt, message, err := ws.ReadMessage()
	fmt.Println(string(message))
	if err != nil {
		panic(err)
	}
	//写入ws数据
	err = ws.WriteMessage(mt, []byte("conn ok"))
	if err != nil {
		panic(err)
	}

	 */

	return nil
}

func (this *WebsocketController) SendAll(ctx *gin.Context) string {
	higo.WebsocketContainer.SendAll(ctx.Query("msg"))
	return "ok"
}
