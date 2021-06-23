package main

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Beans"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/router"
	"github.com/dengpju/higo-utils/utils"
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

	higo.Init(utils.NewSliceString(".", "test", "")).
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

func tt()  {
	ra := make(map[string]bool)
	rb := make(map[string]bool)
	for i := 0; i < 20; i++ {
		if i % 5 == 0 {
			fmt.Println()
		}
	begin_a:
		decade1 := utils.NewRandom().BetweenInt(2, 11) //十位
		unit1 := utils.NewRandom().BetweenInt(0, 8)    //个位
		a := fmt.Sprintf("%d%d", decade1, unit1)
		if _, ok := ra[a]; ok {
			goto begin_a
		}
		ra[a] = true
	begin_b:
		decade2 := utils.NewRandom().BetweenInt(0, decade1-1) //十位
		unit2 := utils.NewRandom().BetweenInt(unit1+1, 9)     //个位
		b := fmt.Sprintf("%d%d", decade2, unit2)
		if _, ok := rb[b]; ok {
			goto begin_b
		}
		rb[b] = true
		fmt.Printf("%s - %s\n", a, b)
	}
}