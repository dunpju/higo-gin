package Middlewares

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-utils/utils/runtimeutil"
	"github.com/gin-gonic/gin"
	"strconv"
)

// RunLog 运行日志
type RunLog struct{}

// NewRunLog 构造函数
func NewRunLog() *RunLog {
	return &RunLog{}
}

func (this *RunLog) Middle(hg *higo.Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		tt := cxt.Query("tt")
		goid, _ := runtimeutil.GoroutineID()
		if route, ok := hg.GetRoute(cxt.Request.Method, cxt.Request.URL.Path); ok {
			fmt.Printf("RunLog:%s\n", route.Desc()+"-"+strconv.FormatUint(goid, 10)+"-"+tt)
		}
		cxt.Next()
	}
}
