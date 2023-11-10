package higo

import (
	"github.com/dunpju/higo-gin/higo/templates"
	"github.com/dunpju/higo-ioc/injector"
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

// AddContainer 注册到Di容器
func AddContainer(builds ...DepBuild) {
	for _, build := range builds {
		cl := build()
		key := container.key(cl)
		if _, ok := container.get(key); !ok {
			container.set(key, build)
		}
	}
}

// Di 获取依赖
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
