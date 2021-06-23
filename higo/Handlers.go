package higo

import "github.com/gin-gonic/gin"

func handleConvert(handler interface{}) interface{} {
	if handle, ok := handler.(func(*gin.Context)); ok {
		return handle
	} else if handle, ok := handler.(func() string); ok {
		return func(ctx *gin.Context) {
			ctx.String(200, handle())
		}
	}
	return nil
}
