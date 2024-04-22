package Codes

import "gitee.com/dengpju/higo-code/code"

// CodeErrorCode error_code码
type CodeErrorCode int64

func (this CodeErrorCode) Message(variables ...interface{}) string {
	return code.Get(this, variables...)
}

func (this CodeErrorCode) Register() *code.Message {
	autoload()
	return code.Container()
}
