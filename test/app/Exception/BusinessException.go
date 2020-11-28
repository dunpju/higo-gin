package Exception

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-utils/utils"
)

type BusinessException struct {
	higo.ServerException // 继承
}

func NewBusinessException(code int, msg string, data ...interface{}) {
	new(BusinessException).Exception(msg, code, utils.Ifindex(data, 0))
}