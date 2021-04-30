package Controllers

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Entity"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-router/router"
	"github.com/gin-gonic/gin"
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
func (this *WebsocketController) Conn(ctx *gin.Context) higo.WsWriteMessage {
	fmt.Println("控制器 Conn")
	fmt.Println("控制器 Conn", ctx.Request.URL.Path)
	loginEntity := Entity.NewLoginEntity()
	err := ctx.ShouldBind(loginEntity)
	if err != nil {
		panic(err)
	}
	fmt.Println("Conn", loginEntity)

	return higo.WsRespStruct(loginEntity)
}

func (this *WebsocketController) Echo(ctx *gin.Context) higo.WsWriteMessage {

	return higo.WsRespString("echo")
}

func (this *WebsocketController) SendAll(ctx *gin.Context) string {
	higo.WsContainer.SendAll(ctx.Query("msg"))
	return "ok"
}
