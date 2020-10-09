package higo

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo/utils"
	"github.com/gin-gonic/gin"
	"runtime"
)

type ServerException struct {

}

func (this *ServerException) Exception(message interface{}, code int, data ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	msg := ErrorToString(message)
	Logrus.Info(fmt.Sprintf("%s (code: %d) at %s:%d", msg, code, file, line))
	panic(gin.H{
		"code": code,
		"msg":  msg,
		"data": utils.Ifindex(data, 0),
	})
}

// 抛出异常
func Throw(message interface{}, code int, data ...interface{}) {
	new(ServerException).Exception(message, code, utils.Ifindex(data, 0))
}

// recover 转 string
func ErrorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	case []uint8:
		return B2S(r.([]uint8))
	default:
		return r.(string)
	}
}

// []uint8 转 string
func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
