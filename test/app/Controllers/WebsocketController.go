package Controllers

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Consts"
	"github.com/dengpju/higo-gin/test/app/Entity"
	"github.com/dengpju/higo-throw/exception"
	"github.com/gin-gonic/gin"
)

type WebsocketController struct {
	*higo.Orm `inject:"Bean.NewOrm()"`
	Redis      *higo.RedisAdapter `inject:"Bean.NewRedisAdapter()"`
}

func NewWebsocketController() *WebsocketController {
	return &WebsocketController{}
}

func (this *WebsocketController) New() higo.IClass {
	return NewWebsocketController()
}

func (this *WebsocketController) Route(hg *higo.Higo) {
	hg.Ws("/conn", this.Conn, hg.Desc("conn"))
	hg.Ws("/echo", this.Echo, hg.Flag("WebsocketController.Echo"), hg.Desc("Echo"))
	hg.Ws("/send_all", this.SendAll, hg.Flag("WebsocketController.SendAll"), hg.Desc("SendAll"))
}

//webSocket请求
func (this *WebsocketController) Conn(ctx *gin.Context) higo.WsWriteMessage {
	fmt.Println("控制器 Conn")
	fmt.Println("控制器 Conn", this)
	fmt.Println("控制器 Conn", ctx.Request.URL.Path)
	//测试异常抛出
	exception.Throw(exception.Message("111"), exception.Code(Consts.CodeError), exception.Data("hhh"))
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
