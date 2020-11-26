package main

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/router"
	"github.com/dengpju/higo-ioc/config"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-ioc/test/services"
)

func main()  {
	serviceConfig:=config.NewServiceConfig()
	injector.BeanFactory.Config(serviceConfig)
	userService:=services.NewUserService()
	injector.BeanFactory.Apply(userService)
	fmt.Println(userService.Order)
	adminService:=services.NewAdminService()
	injector.BeanFactory.Apply(adminService)
	fmt.Println(adminService.Order)
	fmt.Println(adminService.Order.Db)

	higo.Init().
		Middleware(Middlewares.NewAuth(), Middlewares.NewRunLog()).
		SetRoot(".\\test\\").
		//HttpServe("HTTP_HOST", router.NewHttp()).
		HttpsServe("HTTPS_HOST", router.NewHttps()).
		IsAutoGenerateSsl(true).
		Beans(higo.NewHgController(),V3.NewDemoController()).
		Boot()
}
