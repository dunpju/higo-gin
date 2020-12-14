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

// 外部获取配置值
func ConfigValue(key string) string {
	return config.Get(key).(string)
}

// 外部获取配置
func Config(key string) Configure {
	return config.Get(key).(Configure)
}

// 外部获取所有配置
func Configs() Configure {
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

// 获取值
func (this Configure) Value (key string) string {
	return this.Get(key).(string)
}

// 第一个元素
func (this Configure) First () string {
	var first string
	for _,v := range this {
		first = v.(string)
		break
	}
	return first
}

// 获取Configure对象
func (this Configure) Configure (key string) Configure {
	return this.Get(key).(Configure)
}