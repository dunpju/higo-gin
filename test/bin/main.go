package main

import (
	"fmt"
	"gitee.com/dengpju/higo-code/code"
	"github.com/dunpju/higo-config/config"
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-gin/test/app/Beans"
	"github.com/dunpju/higo-gin/test/app/Middlewares"
	"github.com/dunpju/higo-gin/test/router"
	"github.com/dunpju/higo-logger/logger"
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-throw/exception"
	"github.com/dunpju/higo-utils/utils/maputil"
	"github.com/dunpju/higo-utils/utils/randomutil"
	"github.com/dunpju/higo-utils/utils/runtimeutil"
	"github.com/dunpju/higo-utils/utils/sliceutil"
	"github.com/dunpju/higo-wsock/wsock"
	term "github.com/nsf/termbox-go"
	"os/exec"
)

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (r *Resp) SetCode(code int) {
	r.Code = code
}

func (r *Resp) SetMessage(msg string) {
	r.Message = msg
}

func (r *Resp) SetData(data interface{}) {
	r.Data = data
}

func NewResp(code int, message string, data interface{}) higo.IResult {
	return &Resp{Code: code, Message: message, Data: data}
}

func main() {
	higo.NewResult = NewResp
	higo.ResponserTest()("ttt", 0, "hhh")

	//tt();
	checkStatement := fmt.Sprintf("netstat -ano | grep %d", 6123)
	output, _ := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	fmt.Printf("%s", output)

	beanConfig := Beans.NewMyBean()

	//injector.BeanFactory.Config(beanConfig)
	//demoController := V3.NewDemoController()
	//injector.BeanFactory.Apply(demoController)
	//fmt.Println(demoController.DB)

	//higo.WsPitpatSleep = time.Second * 5
	wsock.WsRecoverHandle = func(conn *wsock.WebsocketConn, r interface{}) (respMsg string) {
		goid, _ := runtimeutil.GoroutineID()
		logger.LoggerStack(r, goid)
		if msg, ok := r.(*code.CodeMessage); ok {
			respMsg = maputil.Array().
				Put("code", msg.Code).
				Put("message", msg.Message).
				Put("data", nil).
				String()
		} else if arrayMap, ok := r.(maputil.ArrayMap); ok {
			respMsg = arrayMap.String()
		} else if arrayMap, ok := r.(*maputil.ArrayMap); ok {
			respMsg = arrayMap.String()
		} else {
			fmt.Printf("%T\n", r)
			respMsg = maputil.Array().
				Put("code", 0).
				Put("message", exception.ErrorToString(r)).
				Put("data", nil).
				String()
		}
		return
	}

	higo.Init(sliceutil.NewSliceString(".", "")).
		Middleware(Middlewares.NewRunLog()).
		Middleware(Middlewares.NewWebsocket()).
		AddServe(router.NewHttp(), Middlewares.NewHttp()).
		AddServe(router.NewHttps(), beanConfig).
		//AddServe(router.NewWebsocket()).
		IsAutoTLS(true).
		IsRedisPool().
		Beans(beanConfig).
		//Cron("0/3 * * * * *", func() {
		//	log.Println("3秒执行一次")
		//}).
		Event(higo.AfterLoadConfigure, func(hg *higo.Higo) {
			fmt.Println("测试事件")
			confDefault := config.Db("DB.Default").(*config.Configure)

			_, err := him.DbConfig(him.DefaultConnect).
				SetHost(confDefault.Get("Host").(string)).
				SetPort(confDefault.Get("Port").(string)).
				SetDatabase(confDefault.Get("Database").(string)).
				SetUsername(confDefault.Get("Username").(string)).
				SetPassword(confDefault.Get("Password").(string)).
				SetCharset(confDefault.Get("Charset").(string)).
				SetDriver(confDefault.Get("Driver").(string)).
				SetPrefix(confDefault.Get("Prefix").(string)).
				SetMaxIdle(confDefault.Get("MaxIdle").(int)).
				SetMaxOpen(confDefault.Get("MaxOpen").(int)).
				SetMaxLifetime(confDefault.Get("MaxLifetime").(int)).
				SetLogMode("Info").
				SetSlowThreshold(1).
				SetColorful(true).
				Init()
			if err != nil {
				panic(err)
			}
		}).
		Boot()

}

func tt() {
	ra := make(map[string]bool)
	rb := make(map[string]bool)
	for i := 0; i < 20; i++ {
		if i%5 == 0 {
			fmt.Println()
		}
	begin_a:
		decade1 := randomutil.NewRandom().BetweenInt(2, 11) //十位
		unit1 := randomutil.NewRandom().BetweenInt(0, 8)    //个位
		a := fmt.Sprintf("%d%d", decade1, unit1)
		if _, ok := ra[a]; ok {
			goto begin_a
		}
		ra[a] = true
	begin_b:
		decade2 := randomutil.NewRandom().BetweenInt(0, decade1-1) //十位
		unit2 := randomutil.NewRandom().BetweenInt(unit1+1, 9)     //个位
		b := fmt.Sprintf("%d%d", decade2, unit2)
		if _, ok := rb[b]; ok {
			goto begin_b
		}
		rb[b] = true
		fmt.Printf("%s - %s\n", a, b)
	}
}

func reset() {
	term.Sync() // cosmestic purpose
}

func ter() {
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()
	fmt.Println("Enter any key to see their ASCII code or press ESC button to quit")
keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyF1:
				reset()
				fmt.Println("F1 pressed")
			case term.KeyF2:
				reset()
				fmt.Println("F2 pressed")
			case term.KeyF3:
				reset()
				fmt.Println("F3 pressed")
			case term.KeyF4:
				reset()
				fmt.Println("F4 pressed")
			case term.KeyF5:
				reset()
				fmt.Println("F5 pressed")
			case term.KeyF6:
				reset()
				fmt.Println("F6 pressed")
			case term.KeyF7:
				reset()
				fmt.Println("F7 pressed")
			case term.KeyF8:
				reset()
				fmt.Println("F8 pressed")
			case term.KeyF9:
				reset()
				fmt.Println("F9 pressed")
			case term.KeyF10:
				reset()
				fmt.Println("F10 pressed")
			case term.KeyF11:
				reset()
				fmt.Println("F11 pressed")
			case term.KeyF12:
				reset()
				fmt.Println("F12 pressed")
			case term.KeyInsert:
				reset()
				fmt.Println("Insert pressed")
			case term.KeyDelete:
				reset()
				fmt.Println("Delete pressed")
			case term.KeyHome:
				reset()
				fmt.Println("Home pressed")
			case term.KeyEnd:
				reset()
				fmt.Println("End pressed")
			case term.KeyPgup:
				reset()
				fmt.Println("Page Up pressed")
			case term.KeyPgdn:
				reset()
				fmt.Println("Page Down pressed")
			case term.KeyArrowUp:
				reset()
				fmt.Println("Arrow Up pressed")
			case term.KeyArrowDown:
				reset()
				fmt.Println("Arrow Down pressed")
			case term.KeyArrowLeft:
				reset()
				fmt.Println("Arrow Left pressed")
			case term.KeyArrowRight:
				reset()
				fmt.Println("Arrow Right pressed")
			case term.KeySpace:
				reset()
				fmt.Println("Space pressed")
			case term.KeyBackspace:
				reset()
				fmt.Println("Backspace pressed")
			case term.KeyEnter:
				reset()
				fmt.Println("Enter pressed")
			case term.KeyTab:
				reset()
				fmt.Println("Tab pressed")

			default:
				// we only want to read a single character or one key pressed event
				reset()
				fmt.Println("ASCII : ", ev.Ch)

			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}
