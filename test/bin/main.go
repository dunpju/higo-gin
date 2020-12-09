package main

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Config"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/router"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-utils/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
)

func main()  {
	scanFiles := utils.Dir("./test/app/Controllers").Suffix("go").Scan().List()
	fmt.Println(scanFiles)

	// 通过解析src来创建AST。
	fset := token.NewFileSet() // 职位相对于fset
	f, err := parser.ParseFile(fset, "./test/app/Controllers\\V3\\DemoController.go", nil, 0)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(x ast.Node) bool {
		ts, ok := x.(*ast.TypeSpec)
		if !ok || ts.Type == nil {
			return true
		}

		// 获取结构体名称
		structName := ts.Name.Name
		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		fmt.Println(structName)
		fmt.Println(s)
		tank := fmt.Sprintf("New%s",structName)
		fmt.Println(tank)
		return false
	})


	// 打印AST。
	//_ = ast.Print(fset, f)


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
