package errcode

import (
	"gitee.com/dengpju/higo-code/code"
	"github.com/dunpju/higo-gin/higo/errcode"
)

func init() {
	errcode.Autoload = autoload
}

const (
	AuthError      errcode.ErrorCode = iota + 1 //认证失败
	EnumError                                   //枚举错误
	ParamError                                  //参数错误
	PrimaryIdError                              //主键id错误
	UniqueError                                 //重复
	NotExistError                               //不存在
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
