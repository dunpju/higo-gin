package higo

var container Dependency

type Dependency map[string]IClass

func init() {
	container = make(Dependency)
}

// 注册到Di容器
func AddDiToContainer(class IClass)  {
	rt, _ := class.Reflection()
	typ := rt.Name()
	container[typ] = class
}

// 获取依赖
func Di(name string) IClass {
	return container[name]
}