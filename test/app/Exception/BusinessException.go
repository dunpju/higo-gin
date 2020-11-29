package Exception

import (
	"github.com/dengpju/higo-throw/throw"
	"github.com/dengpju/higo-utils/utils"
)

type BusinessException struct {
	throw.ServerException // 继承
}

func NewBusinessException(code int, msg string, data ...interface{}) {
	new(BusinessException).Exception(msg, code, utils.Ifindex(data, 0))
}