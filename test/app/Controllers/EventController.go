package Controllers

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-gin/higo/event"
	"github.com/dunpju/higo-gin/higo/request"
	"github.com/dunpju/higo-gin/higo/responser"
	"github.com/dunpju/higo-gin/test/app/Services"
	"github.com/dunpju/higo-utils/utils/runtimeutil"
	"github.com/gin-gonic/gin"
	"time"
)

type EventController struct {
	T string
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
		hg.Get("/test1", this.Test1, hg.Flag("EventController.Test"), hg.Desc("事件测试1"))
		hg.Get("/test2", this.Test2, hg.Flag("EventController.Test"), hg.Desc("事件测试2"))
		hg.Get("/test3", this.Test3, hg.Flag("EventController.Test"), hg.Desc("事件测试3"))
		hg.Get("/test4", Test4, hg.Flag("EventController.Test"), hg.Desc("事件测试4"))
	})
}

func (this *EventController) Test1() string {
	this.T = "test1"
	fmt.Println(this)
	return "Test1"
}

func (this *EventController) Test2() interface{} {
	fmt.Println(runtimeutil.GoroutineID())
	ctx := request.Context()
	tt := ctx.Query("tt")
	fmt.Println(this)
	return tt
}

var i = 0

func (this *EventController) Test3(ctx1 *gin.Context) {
	fmt.Println("请求数量")
	gid, _ := runtimeutil.GoroutineID()
	tid, err := runtimeutil.ThreadID()
	fmt.Println(err)
	if gid%2 == 0 {
		fmt.Printf("线程: %d 协成: %d  %s\n", tid, gid, "休眠")
		time.Sleep(2 * time.Second)
		i++
		if i == 1 {
			panic("测试异常")
		}
	}
	//fmt.Println(len(higo.Request))
	ctx := request.Context()
	tt := ctx.Query("tt")

	fmt.Printf("线程: %d 协成: %d  %s\n", tid, gid, tt)
	go func() {
		gid, _ := runtimeutil.GoroutineID()
		tid, _ := runtimeutil.ThreadID()
		fmt.Printf("线程: %d 子协成: %d 数据:%s\n", tid, gid, tt)
	}()
	//exception.Throw(exception.Message(tt), exception.Code(1))
	responser.SuccessJson("success", 10000, tt)
}

func Test4() {
	ctx := request.Context()
	fmt.Println(runtimeutil.GoroutineID())
	tt := ctx.Query("tt")
	fmt.Println(tt)
	//exception.Throw(exception.Message(tt), exception.Code(1))
	responser.SuccessJson("success", 10000, tt)
}

// 订阅数据
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
	/**
	go func() {
		//发布
		ev.Pub("info", nil)
	}()
	ch := ev.Sub("info", getUserInfo)//订阅
	*/
	ch := Services.GetDemoListCh()
	Services.Bus.Pub(Services.GetDemoList, ch)
	defer Services.Bus.UnSub(Services.GetDemoList, ch)
	higo.Responser(ctx).SuccessJson("success", 10000, ch.Data(time.Second*1))
}

func testPub() interface{} {
	time.Sleep(time.Second * 5)
	return "商品列表"
}

func getInfo() string {
	type example_model struct {
		Id   int
		Name string
	}
	var models []*example_model
	models = append(models, &example_model{Id: 1, Name: "foo"}, &example_model{Id: 1, Name: "bar"})
	return "这是信息"
}

// 分体
func getUserInfo(id int) interface{} {
	return gin.H{"id": id, "商品分体": "ffff"}
}
