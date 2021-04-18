package main

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Config"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/router"
	"os/exec"
)

func main() {

	checkStatement := fmt.Sprintf("netstat -ano | grep %d", 6123)
	output, _ := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	fmt.Printf("%s", output)

	beanConfig := Config.NewBean()

	//injector.BeanFactory.Config(beanConfig)
	//demoController := V3.NewDemoController()
	//injector.BeanFactory.Apply(demoController)
	//fmt.Println(demoController.DB)

	higo.Init().
		Middleware(Middlewares.NewAuth(), Middlewares.NewRunLog()).
		LoadConfigur(".\\test\\").
		AddServe(router.NewHttp()).
		AddServe(router.NewHttps()).
		IsAutoTLS(true).
		IsRedisPool().
		Beans(beanConfig).
		//Cron("0/3 * * * * *", func() {
		//	log.Println("3秒执行一次")
		//}).
		Boot()

}
