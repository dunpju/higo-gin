package main

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/router"
)

func main()  {
	higo.Init().
		SetRoot(".\\test\\").
		HttpsServe("HTTPS_HOST", router.NewHttps()).
		IsAutoGenerateSsl(true).
		Boot()
}
