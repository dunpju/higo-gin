package Controllers

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/higo/event"
	"github.com/gin-gonic/gin"
	"time"
)

type EventController struct {
}

func NewEventController() *EventController {
	return &EventController{}
}

func (this *EventController) New() higo.IClass {
	return NewEventController()
}

func (this *EventController) Route(hg *higo.Higo) {
	hg.AddGroup("/event", func() {
		hg.Get("/test", this.Test, hg.Flag("EventController.Test"), hg.Desc("事件测试"))
	})
}

//订阅数据
var ev = event.NewEventBus() //需要全局

func (this *EventController) Test(ctx *gin.Context) {

	/**
	ch := ev.Sub("user")
	go func() {
		//发布
		ev.Pub("user", testPub())
	}()

	higo.Responser(ctx).SuccessJson("success", 10000, ch.Data(time.Second*1))

	*/

	go func() {
		//发布
		ev.Pub("info", nil)
	}()
	ch := ev.Sub("info", getUserInfo)//订阅
	higo.Responser(ctx).SuccessJson("success", 10000, ch.Data(time.Second*1))
}

func testPub() interface{} {
	time.Sleep(time.Second * 5)
	return "商品列表"
}

func getInfo() string {
	return "这是信息"
}

//分体
func getUserInfo(id int) interface{} {
	return gin.H{"id": id, "商品分体": "ffff"}
}
