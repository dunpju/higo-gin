package main

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Beans"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/router"
	"github.com/dengpju/higo-utils/utils/randomutil"
	"github.com/dengpju/higo-utils/utils/sliceutil"
	term "github.com/nsf/termbox-go"
	"os/exec"
)

func main() {
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

	higo.Init(sliceutil.NewSliceString(".", "test", "")).
		Middleware(Middlewares.NewRunLog()).
		AddServe(router.NewHttp(), Middlewares.NewHttp()).
		AddServe(router.NewHttps(), beanConfig).
		AddServe(router.NewWebsocket()).
		IsAutoTLS(true).
		IsRedisPool().
		Beans(beanConfig).
		//Cron("0/3 * * * * *", func() {
		//	log.Println("3秒执行一次")
		//}).
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