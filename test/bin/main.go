package main

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/router"
)

func main()  {

	higo.Init().
		Middleware(Middlewares.NewAuth(), Middlewares.NewRunLog()).
		SetRoot(".\\test\\").
		HttpServe("HTTP_HOST", router.NewHttp()).
		HttpsServe("HTTPS_HOST", router.NewHttps()).
		IsAutoGenerateSsl(true).
		Beans(V3.NewDemoController()).
		Boot()
}
