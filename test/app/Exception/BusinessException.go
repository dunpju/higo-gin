package Exception

import (
	"gitee.com/dengpju/higo-parameter/parameter"
	"github.com/dengpju/higo-throw/throw"
)

type Business struct {
	throw.Throwable // 继承
}

func BusinessException(p ...*parameter.Parameter) {
	new(Business).Exception(p...)
}