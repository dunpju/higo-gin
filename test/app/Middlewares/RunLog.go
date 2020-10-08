package Middlewares

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

// 运行日志
type RunLog struct {}

// 构造函数
func NewRunLog() *RunLog {
	return &RunLog{}
}

func (this *RunLog) Loader(hg *higo.Higo) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(higo.Container().Route(c.Request.URL.Path).Desc)
		c.Next()
	}
}
