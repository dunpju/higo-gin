package higo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

// 抛出异常
func Throw(message interface{}, code int, data ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	msg := ErrorToString(message)
	Logrus.Info(fmt.Sprintf("%s (code: %d) at %s:%d", msg, code, file, line))
	panic(gin.H{
		"code": code,
		"msg":  msg,
		"data": data[0],
	})
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
