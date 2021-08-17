package Consts

import "github.com/dengpju/higo-enum/enum"

type SuccessConst int

const (
	Success SuccessConst = iota + 20000
)

func (this SuccessConst) String() string {
	switch this {
	case Success:
		return "成功"
	}
	return "未定义"
}

func (this SuccessConst) Message() *enum.CodeDoc {
	return enum.New(this)
}

type DemoConst int

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

func (this DemoConst) String() string {
	switch this {
	case ServerError:
		return "系统错误"
	case AuthError:
		return "认证错误"
	case UnknownError:
		return "未知错误"
	case RsaError:
		return "解密错误"
	case ParameterError:
		return "参数错误"
	case NotFound:
		return "未找到"
	case CodeError:
		return "失败"
	case InvalidToken:
		return "无效token"
	case InvalidApi:
		return "无效api"
	case InvalidMap:
		return "无效api映射"
	case TestError:
		return "测试异常"
	}
	return "未定义"
}

func (this DemoConst) Message() *enum.CodeDoc {
	return enum.New(this)
}
