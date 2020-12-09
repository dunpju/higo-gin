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
	"reflect"
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
			r := reflect.TypeOf(x)
			rv := reflect.ValueOf(x)
			fmt.Println(r.Elem().Name())
			fmt.Println(rv.Interface())
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}
		return true
	})


	/**
	ast.Inspect(f, func(x ast.Node) bool {
		s, ok := x.(*ast.StructType)
		if !ok {
			return true
		}

		rs := reflect.ValueOf(x)
		ir := reflect.TypeOf(rs.Interface())
		fmt.Println(ir.String())

		for _, field := range s.Fields.List {
			fmt.Printf("Field: %s\n", field.Names)
			//fmt.Printf("Tag:   %s\n", field.Doc)
		}
		return false
	})
	 */

	// 打印AST。
	_ = ast.Print(fset, f)

	/**
	var buf bytes.Buffer

	_ = ast.Fprint(&buf, fset, f, ast.NotNilFilter)
	// 删除包围函数体的大括号{}，unindent，
	// 并修剪前导和尾随空白区域。
	s := buf.String()
	s = s[1 : len(s)-1]
	s = strings.TrimSpace(strings.Replace(s, "\n\t", "\n", -1))

	// 将清理后的正文文本打印到标准输出。
	//fmt.Println(s)

	 */


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
