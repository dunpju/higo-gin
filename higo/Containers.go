package higo

import "github.com/dengpju/higo-throw/throw"

// 容器
type Containers struct {
	// 配置
	C map[interface{}]interface{}
	// 路由
	R map[string]Route
	// Di
	Di map[string]IBuilder
}

// 构造函数
func NewContainer() *Containers {
	return &Containers{
		C:  make(map[interface{}]interface{}),
		R:  make(map[string]Route),
		Di: make(map[string]IBuilder),
	}
}

// 获取所有配置
func (this *Containers) Configure() map[interface{}]interface{} {
	return this.C
}

// 获取配置
func (this *Containers) Config(key string) map[interface{}]interface{} {
	v, ok := this.C[key]
	if !ok {
		throw.Throw("获取"+key+"配置失败", 0)
	}
	return v.(map[interface{}]interface{})
}

// 添加路由容器
func (this *Containers) AddRoutes(relativePath string, route Route) *Containers {
	this.R[relativePath] = route
	return this
}

// 获取所有路由
func (this *Containers) Routes() map[string]Route {
	return this.R
}

// 获取路由
func (this *Containers) Route(relativePath string) Route {
	route, ok := this.R[relativePath]
	if !ok {
		throw.Throw(relativePath+"未定义", 0)
	}
	return route
}

// 获取依赖
func Di(name string) IBuilder {
	return Container().Di[name]
}