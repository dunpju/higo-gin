package errcode

import (
	"gitee.com/dengpju/higo-code/code"
)

// 自动加载
var Autoload func()

// 错误码
type ErrorCode int64

func (this ErrorCode) Message(variables ...interface{}) string {
	return code.Get(this, variables...)
}

func (this ErrorCode) Register() code.Message {
	Autoload()
	return code.Container()
}
