package main

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Config"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/router"
	"github.com/dengpju/higo-ioc/injector"
)

func main()  {
	beanConfig := Config.NewBean()
	injector.BeanFactory.Config(beanConfig)
	demoController := V3.NewDemoController()
	injector.BeanFactory.Apply(demoController)
	fmt.Println(demoController.DemoService)

	higo.Init().
		Middleware(Middlewares.NewAuth(), Middlewares.NewRunLog()).
		SetRoot(".\\test\\").
		//HttpServe("HTTP_HOST", router.NewHttp()).
		HttpsServe("HTTPS_HOST", router.NewHttps()).
		IsAutoGenerateSsl(true).
		Beans(beanConfig).
		Boot()
}
