package request

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

func Context() *gin.Context {
	return higo.Request.Context()
}

func Set(context *gin.Context) {
	higo.Request.Set(context)
}
