package Middlewares

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-utils/utils/runtimeutil"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 运行日志
type RunLog struct{}

// 构造函数
func NewRunLog() *RunLog {
	return &RunLog{}
}

func (this *RunLog) Middle(hg *higo.Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		tt := cxt.Query("tt")
		fmt.Printf("RunLog:%s\n",
			higo.RouterContainer.Get(cxt.Request.Method, cxt.Request.URL.Path).Desc()+"-"+strconv.FormatUint(runtimeutil.GoroutineID(), 10)+"-"+tt)
		cxt.Next()
	}
}
