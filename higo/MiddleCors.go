package higo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 跨域
type Cors struct{}

// 构造函数
func NewCors() *Cors {
	return &Cors{}
}

func (this *Cors) Middle(hg *Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		MiddleCorsFunc(hg)(cxt)
		cxt.Next()
	}
}

func middleCorsFunc(hg *Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		method := cxt.Request.Method
		origin := cxt.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			cxt.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			cxt.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			cxt.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			cxt.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			cxt.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			cxt.AbortWithStatus(http.StatusNoContent)
		}
	}
}
