package higo

import (
	"fmt"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-utils/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
)

var (
	container       Dependency
	RefDepBuildType reflect.Type
	depb            DepBuild
)

func init() {
	RefDepBuildType = reflect.TypeOf(depb)
}

type DepBuild func() IClass

type Dependency map[string]DepBuild

func Scan() {
	scanFiles := utils.Dir("./test/app/Controllers").Suffix("go").Scan().Get()
	fmt.Println(scanFiles)
	fmt.Println(container)

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
		structType, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		fmt.Println(structName)
		fmt.Println("42", structType)
		fmt.Printf("%T\n", structType)
		tank := fmt.Sprintf("New%s", structName)
		fmt.Println(tank)
		//reflect.New(structName).Elem().Interface()
		return false
	})

	// 打印AST。
	//_ = ast.Print(fset, f)
}

// 注册到Di容器
func AddContainer(builds ...DepBuild) {
	for _, build := range builds {
		class := build()
		v := reflect.ValueOf(class)
		if _, ok := container[v.Type().String()]; !ok {
			container[v.Type().String()] = build
		}
	}
}

// 获取依赖
func Di(name string) IClass {
	class := container[name]()
	injector.BeanFactory.Apply(class)
	injector.BeanFactory.Set(class)
	return class
}
