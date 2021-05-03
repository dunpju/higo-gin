package higo

import (
	"github.com/gin-gonic/gin"
)

// 跨域
type Cors struct{}

// 构造函数
func NewCors() *Cors {
	return &Cors{}
}

func (this *Cors) Middle(hg *Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		MiddleCorsFunc(cxt)
		cxt.Next()
	}
}
