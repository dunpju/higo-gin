package Exception

import "github.com/dengpju/higo-gin/higo"

type BusinessException struct {

}

func NewBusinessException(code int, msg string) {
	higo.Throw(msg, code)
}

// 业务异常
func (this *BusinessException) Throw(message string, code int) {
	higo.Throw(message, code)
}