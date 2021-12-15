package errcodg

import (
	"gitee.com/dengpju/higo-code/code"
)

//错误码
type ErrorCode int64

func (this ErrorCode) Message(variables ...interface{}) string {
	return code.Get(this, variables...)
}

func (this ErrorCode) Register() code.Message {
	autoload()
	return code.Container()
}

const (
	AuthError      ErrorCode = iota + 1 //认证失败
	EnumError                           //枚举错误
	ParamError                          //参数错误
	PrimaryIdError                      //主键id错误
	UniqueError                         //重复
	NotExistError                       //不存在

)

func autoload() {
	code.Container().
		Put(AuthError, "认证失败").
		Put(EnumError, "枚举错误").
		Put(ParamError, "参数错误").
		Put(PrimaryIdError, "主键id错误").
		Put(UniqueError, "重复").
		Put(NotExistError, "不存在")
}
