package Controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Entity"
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
	router.Ws("/conn", this.Conn, router.Flag("WebsocketController.Conn"), router.Desc("conn"))
	router.Ws("/echo", this.Echo, router.Flag("WebsocketController.Echo"), router.Desc("Echo"))
	router.Ws("/send_all", this.SendAll, router.Flag("WebsocketController.SendAll"), router.Desc("SendAll"))
	return hg
}

//webSocket请求
func (this *WebsocketController) Conn(ctx *gin.Context) higo.Websocket {
	fmt.Println("控制器 Conn")
	loginEntity := Entity.NewLoginEntity()
	err := ctx.ShouldBind(loginEntity)
	if err != nil {
		panic(err)
	}
	fmt.Println("Conn", loginEntity)

	result,err := json.Marshal(loginEntity)
	return string(result)
}

func (this *WebsocketController) Echo(ctx *gin.Context) higo.Websocket {

	return "echo"
}

func (this *WebsocketController) SendAll(ctx *gin.Context) string {
	higo.WebsocketContainer.SendAll(ctx.Query("msg"))
	return "ok"
}
