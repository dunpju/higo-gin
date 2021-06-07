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

func (this *EventController) Test(ctx *gin.Context) {
	//订阅数据
	ev := event.NewEventBus()
	ch := ev.Sub("user")
	go func() {
		//发布
		ev.Pub("user", testPub())
	}()

	higo.Responser(ctx).SuccessJson("success", 10000, ch.Data(time.Second*1))
}

func testPub() interface{} {
	time.Sleep(time.Second * 5)
	return "商品列表"
}
