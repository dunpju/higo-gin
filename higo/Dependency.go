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

var container Dependency

type Dependency map[string]IClass

func Scan()  {
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
		//reflect.New(structName).Elem().Interface()
		return false
	})


	// 打印AST。
	//_ = ast.Print(fset, f)
}

func Test()  {
	t := "*V3.DemoController"
	self := Di(t).Self()
	i := self
	fmt.Println(i)
}

// 注册到Di容器
func AddContainer(class IClass)  {
	injector.BeanFactory.Apply(class)
	injector.BeanFactory.Set(class)
	v := reflect.ValueOf(class)
	container[v.Type().String()]=class
}

// 获取依赖
func Di(name string) IClass {
	return container[name]
}