package higo

import "github.com/dengpju/higo-throw/throw"

// 容器
type Containers struct {
	// 配置
	Conf map[interface{}]interface{}
	// 路由
	Rout map[string]Route
	// Di
	Di map[string]IClass
}

// 构造函数
func NewContainer() *Containers {
	return &Containers{
		Conf: make(map[interface{}]interface{}),
		Rout: make(map[string]Route),
		Di:   make(map[string]IClass),
	}
}

// 获取所有配置
func (this *Containers) Configure() map[interface{}]interface{} {
	return this.Conf
}

// 获取配置
func (this *Containers) Config(key string) map[interface{}]interface{} {
	v, ok := this.Conf[key]
	if !ok {
		throw.Throw("获取"+key+"配置失败", 0)
	}
	return v.(map[interface{}]interface{})
}

// 添加路由容器
func (this *Containers) AddRoutes(relativePath string, route Route) *Containers {
	this.Rout[relativePath] = route
	return this
}

// 获取所有路由
func (this *Containers) Routes() map[string]Route {
	return this.Rout
}

// 获取路由
func (this *Containers) Route(relativePath string) Route {
	route, ok := this.Rout[relativePath]
	if !ok {
		throw.Throw(relativePath+"未定义", 0)
	}
	return route
}

// 注册到Di容器
func AddDiToContainer(class IClass)  {
	rt, _ := class.Reflection()
	typ := rt.Name()
	Container().Di[typ] = class
}

// 获取依赖
func Di(name string) IClass {
	return Container().Di[name]
}
