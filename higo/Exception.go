package higo

import (
	"gitee.com/dengpju/higo-parameter/parameter"
	"github.com/dunpju/higo-throw/exception"
)

func Throw(parameters ...*parameter.Parameter) {
	exception.Throw(parameters...)
}
