package main

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Config"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/router"
	"github.com/dengpju/higo-ioc/injector"
	"go/importer"
	"os/exec"
)

func main()  {

	pkg, err := importer.Default().Import("exec")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, declName := range pkg.Scope().Names() {
		fmt.Println(declName)
	}

	return

	checkStatement := fmt.Sprintf("netstat -ano | grep %d", 6123)
	output, _ := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	fmt.Printf("%s",output)


	beanConfig := Config.NewBean()

	injector.BeanFactory.Config(beanConfig)
	demoController := V3.NewDemoController()
	injector.BeanFactory.Apply(demoController)
	fmt.Println(demoController.DB)

	higo.Init().
		Middleware(Middlewares.NewAuth(), Middlewares.NewRunLog()).
		SetRoot(".\\test\\").
		//HttpServe("HTTP_HOST", router.NewHttp()).
		HttpsServe("HTTPS_HOST", router.NewHttps()).
		IsAutoGenerateSsl(true).
		Beans(beanConfig).
		Boot()
}
