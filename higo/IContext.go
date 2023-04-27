package higo

import "github.com/gin-gonic/gin"

// 上下文
type IContext interface {
	OnRequest(*gin.Context) error
	OnResponse(result interface{}) (interface{}, error)
}
