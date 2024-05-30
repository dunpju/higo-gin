package errcode

import (
	"fmt"
	"gitee.com/dengpju/higo-code/code"
)

// Autoload 自动加载
var Autoload func()

// ErrorCode 错误码
type ErrorCode int64

func (this ErrorCode) Message(variables ...interface{}) string {
	return code.Get(this, variables...)
}

func (this ErrorCode) Register() *code.Message {
	Autoload()
	return code.Container()
}

func (this ErrorCode) Int() int {
	return int(this)
}

func (this ErrorCode) Int64() int64 {
	return int64(this)
}

func (this ErrorCode) Error(variables ...interface{}) error {
	return fmt.Errorf(this.Message(variables...))
}

func (this ErrorCode) Panic(variables ...interface{}) {
	panic(fmt.Errorf(this.Message(variables...)))
}
