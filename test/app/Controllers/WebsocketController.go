package Controllers

import (
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
	router.Get("/ping", this.Ping, router.Flag("WebsocketController.Ping"), router.Desc("ping"))
	router.Get("/send_all", this.SendAll, router.Flag("WebsocketController.SendAll"), router.Desc("ping"))
	return hg
}

//webSocket请求ping 返回pong
func (this *WebsocketController) Ping(ctx *gin.Context) higo.WebsocketPong {
	return higo.WebsocketPongHandler(ctx)
}

func (this *WebsocketController) SendAll(ctx *gin.Context) higo.WebsocketPong {
	higo.WebsocketClientContainer.SendAll("hello")
	return "ok"
}
