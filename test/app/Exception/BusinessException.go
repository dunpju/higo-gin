package Exception

import "github.com/dengpju/higo-gin/higo"

type BusinessException struct {

}

func NewBusinessException(code int, msg string, data ...interface{}) {
	new(BusinessException).Throw(msg,code,data)
}

// 业务异常
func (this *BusinessException) Throw(message string, code int, data ...interface{}) {
	higo.Throw(message, code, data)
}