package Middlewares

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

// http服务中间件
type Http struct{}

// 构造函数
func NewHttp() *Http {
	return &Http{}
}

func (this *Http) Middle(hg *higo.Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		fmt.Println("http 中间件")
		cxt.Next()
	}
}
