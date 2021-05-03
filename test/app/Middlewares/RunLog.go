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

func (this *RunLog) Middle(hg *higo.Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		fmt.Printf("RunLog:%s\n",higo.RouterContainer.Get(cxt.Request.URL.Path).Desc())
		cxt.Next()
	}
}
