package Controllers

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-gin/test/app/Entity"
	"github.com/dunpju/higo-wsock/wsock"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type WebsocketController struct {
	//*higo.Orm `inject:"Bean.NewOrm()"`
	//Redis      *higo.RedisAdapter `inject:"Bean.NewRedisAdapter()"`
	K string
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

// Conn webSocket请求
// {"username":"3","password":"2","captcha_code":"731","time":"5350"}
func (this *WebsocketController) Conn(ctx *gin.Context) {
	/*var params interface{}
	err := ctx.ShouldBind(&params)
	if err != nil {
		panic(err)
	}
	fmt.Println(params)*/
	//fmt.Println(ctx.GetRawData())
	/*wsock.Response(ctx).WriteMessage("dddd: " + fmt.Sprintf("%v", params))*/
	fmt.Println("控制器 Conn")
	fmt.Println("控制器 Conn", this)
	fmt.Println("控制器 Conn", ctx.Request.URL.Path)
	//测试异常抛出
	//exception.Throw(exception.Message("111"), exception.Code(Consts.CodeError), exception.Data("hhh"))
	loginEntity := Entity.NewLoginEntity()
	err := ctx.ShouldBindBodyWith(loginEntity, binding.JSON)
	if err != nil {
		panic(err)
	}
	fmt.Println("Conn", loginEntity)
}

func (this *WebsocketController) Echo(ctx *gin.Context) {
	wsock.Response(ctx).WriteMessage("dddd")
}

func (this *WebsocketController) SendAll(ctx *gin.Context) {
	wsock.WsContainer.SendAll(ctx.Query("msg"))
}
