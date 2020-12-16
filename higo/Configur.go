package higo

import (
	"github.com/dengpju/higo-throw/throw"
)

var config Configure

type Configure map[interface{}]interface{}

func NewConfigure() *Configure {
	return &config
}

// 外部获取配置值
func ConfigValue(key string) string {
	configure := config.Get(key)
	if nil == configure {
		return ""
	}
	return configure.(string)
}

// 外部获取配置
func Config(key string) Configure {
	configure := config.Get(key)
	if nil == configure {
		return nil
	}
	return configure.(Configure)
}

// 外部获取所有配置
func ConfigAll() Configure {
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
func (this Configure) StrValue (key string) string {
	return this.Get(key).(string)
}

func (this Configure) IntValue (key string) int {
	return this.Get(key).(int)
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