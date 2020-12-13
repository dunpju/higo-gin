package higo

// 容器
type Containers struct {
	// Di
	Di map[string]IClass
}

//type di

// 构造函数
func NewContainer() *Containers {
	return &Containers{
		Di:   make(map[string]IClass),
	}
}

// 注册到Di容器
func AddDiToContainer(class IClass)  {
	//rt, _ := class.Reflection()
	//typ := rt.Name()
	//Container().Di[typ] = class
}

// 获取依赖
//func Di(name string) IClass {
//	return Container().Di[name]
//}
