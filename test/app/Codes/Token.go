package Codes

import "gitee.com/dengpju/higo-code/code"

// token码
type CodeToken int64

func (this CodeToken) Message(variables ...interface{}) string {
	return code.Get(this, variables...)
}

const (
	TokenEmpty CodeToken = iota + 400001 //token为空
)

func (this CodeToken) Register() *code.Message {
	return code.Container().
		Put(TokenEmpty, "token为空")
}
