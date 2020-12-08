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
	"go/importer"
	"go/parser"
	"go/token"
	"os/exec"
)

func main()  {
	scanFiles := utils.Dir("./test/app/Controllers").Suffix("go").Scan().List()
	fmt.Println(scanFiles)
	pkg, err := importer.Default().Import("os/exec")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, declName := range pkg.Scope().Names() {
		fmt.Println(declName)
	}

	// 通过解析src来创建AST。
	fset := token.NewFileSet() // 职位相对于fset
	f, err := parser.ParseFile(fset, "./test/app/Controllers/V3/DemoController.go", nil, 0)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident:
			s = x.Name
		case *ast.StructType:
			s = "struct"
			fmt.Println(x.Fields.)
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}
		return true
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
