package BaseException

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-throw/exception"
)

func Throw(message string, code int) {
	higo.Throw(exception.Code(code), exception.Message(message), exception.Data(""))
}
