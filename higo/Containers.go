package higo

// 容器
type Containers struct {
	// 配置
	Configure map[interface{}]interface{}
	// 路由
	Routes map[string]Route
}

// 构造函数
func NewContainer() *Containers {
	return &Containers{
		Configure: make(map[interface{}]interface{}),
		Routes:    make(map[string]Route),
	}
}

// 获取所有配置
func (this *Containers) GetConfigure() map[interface{}]interface{} {
	return this.Configure
}

// 获取配置
func (this *Containers) Config(key string) map[interface{}]interface{} {
	v, ok := this.Configure[key]
	if !ok {
		Throw("获取" + key + "配置失败",0)
	}
	return v.(map[interface{}]interface{})
}

// 添加路由容器
func (this *Containers) AddRoutes(relativePath string, route Route) *Containers {
	this.Routes[relativePath] = route
	return this
}

// 获取所有路由
func (this *Containers) GetRoutes() map[string]Route {
	return this.Routes
}

// 获取路由
func (this *Containers) GetRoute(relativePath string) Route {
	r, ok := this.Routes[relativePath]
	if !ok {
		Throw(relativePath + "未定义路由", 0)
	}
	return r
}
