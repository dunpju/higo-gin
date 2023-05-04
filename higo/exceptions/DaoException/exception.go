package DaoException

import "github.com/dunpju/higo-gin/higo/exceptions/BaseException"

func Throw(message string, code int) {
	BaseException.Throw(message, code)
}
