package DaoException

import "github.com/dengpju/higo-gin/higo/exceptions/BusinessException"

func Throw(message string, code int) {
	BusinessException.Throw(message, code)
}
