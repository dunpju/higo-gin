package higo

import (
	"github.com/dengpju/higo-throw/throw"
)

var config Configure

type Configure map[interface{}]interface{}

func init() {
	initConfig()
}

// 初始化
func initConfig() {
	Once.Do(func() {
		config = make(Configure)
	})
}

func NewConfigure() *Configure {
	return &config
}

// 外部获取配置
func Config(key string) interface{} {
	return config.Get(key)
}

// 外部获取所有配置
func Configures() Configure {
	return config.All()
}

// 获取所有配置
func (this Configure) All() Configure {
	return this
}

// 获取配置
func (this Configure) Get(key string) interface{} {
	v, ok := this[key]
	if !ok {
		throw.Throw("获取"+key+"配置失败", 0)
	}
	return v
}
