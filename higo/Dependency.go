package higo

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo/templates"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-utils/utils/dirutil"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"sync"
)

var (
	container *Dependency
)

type DepBuild func() IClass

type Dependency struct {
	container *sync.Map
}

func NewDependency() *Dependency {
	return &Dependency{container: &sync.Map{}}
}

func (this *Dependency) set(key string, d DepBuild) {
	this.container.Store(key, d)
}

func (this *Dependency) get(key string) (DepBuild, bool) {
	v, ok := this.container.Load(key)
	if ok {
		return v.(DepBuild), true
	}
	return nil, false
}

func (this *Dependency) key(class interface{}) string {
	v := reflect.ValueOf(class)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Type().PkgPath() + "/" + v.Type().Name()
}

func Scan() {
	scanFiles := dirutil.Dir("./test/app/Controllers").Suffix("go").Scan().Get()
	fmt.Println(scanFiles)
	fmt.Println(container)

	// 通过解析src来创建AST。
	fset := token.NewFileSet() // 相对于fset
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
		cl := build()
		key := container.key(cl)
		if _, ok := container.get(key); !ok {
			container.set(key, build)
		}
	}
}

// 获取依赖
func Di(name string) IClass {
	name = strings.Replace(name, templates.GetModName(), "", 1)
	kk := "/" + strings.TrimLeft(name, "/")
	k := templates.GetModName() + kk
	v, ok := container.get(k)
	if ok {
		class := v()
		injector.BeanFactory.Apply(class)
		injector.BeanFactory.Set(class)
		return class
	}
	return nil
}
