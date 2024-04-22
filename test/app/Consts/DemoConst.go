package Consts

import (
	"gitee.com/dengpju/higo-code/code"
)

type SuccessConst int64

const (
	Success SuccessConst = iota + 20000
)

func (this SuccessConst) Message(variables ...interface{}) string {
	return code.Get(this, variables...)
}

func (this SuccessConst) Register() *code.Message {
	return code.Container().
		Put(Success, "成功")
}

type DemoConst int64

const (
	ServerError DemoConst = iota + 50000
	AuthError
	UnknownError
	RsaError
	ParameterError
	NotFound
	CodeError
	InvalidToken
	InvalidApi
	InvalidMap
	TestError
)

func (this DemoConst) Message(variables ...interface{}) string {
	return code.Get(this, variables...)
}

func (this DemoConst) Register() *code.Message {
	return code.Container().
		Put(ServerError, "系统错误").
		Put(AuthError, "认证错误").
		Put(UnknownError, "未知错误").
		Put(RsaError, "解密错误").
		Put(ParameterError, "参数错误").
		Put(NotFound, "未找到").
		Put(CodeError, "失败").
		Put(InvalidToken, "无效token").
		Put(InvalidApi, "无效api").
		Put(InvalidMap, "无效api映射").
		Put(TestError, "测试异常")
}
