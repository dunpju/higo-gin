package main

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/router"
)

func main()  {

	//d := &V3.DemoController{}
	//higo.NewDirector(d);
	higo.Init().
		Beans(V3.NewDemoController()).
		Middleware(Middlewares.NewAuth(), Middlewares.NewRunLog()).
		SetRoot(".\\test\\").
		HttpServe("HTTP_HOST", router.NewHttp()).
		HttpsServe("HTTPS_HOST", router.NewHttps()).
		IsAutoGenerateSsl(true).
		Boot()
}
