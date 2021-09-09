package BusinessException

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-throw/exception"
)

func Throw(message string, code int) {
	higo.Throw(exception.Code(code), exception.Message(message), exception.Data(""))
}
