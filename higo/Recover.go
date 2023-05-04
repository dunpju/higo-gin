package higo

import (
	"gitee.com/dengpju/higo-code/code"
	"github.com/dunpju/higo-logger/logger"
	"github.com/dunpju/higo-throw/exception"
	"github.com/dunpju/higo-utils/utils/maputil"
	"github.com/dunpju/higo-utils/utils/runtimeutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	//Recover处理函数(可自定义替换)
	RecoverHandle RecoverFunc
	recoverOnce   sync.Once
)

func init() {
	recoverOnce.Do(func() {
		RecoverHandle = func(cxt *gin.Context, r interface{}) {

			//记录debug调用栈
			goid, _ := runtimeutil.GoroutineID()
			logger.LoggerStack(r, goid)

			//封装通用json返回
			if h, ok := r.(gin.H); ok {
				cxt.JSON(http.StatusOK, h)
			} else if cd, ok := r.(*code.CodeMessage); ok {
				cxt.JSON(http.StatusOK, gin.H{
					"code":    cd.Code,
					"message": cd.Message,
					"data":    nil,
				})
			} else if arrayMap, ok := r.(*maputil.ArrayMap); ok {
				cxt.JSON(http.StatusOK, arrayMap.Value())
			} else if validate, ok := r.(*ValidateError); ok {
				cxt.JSON(http.StatusOK, gin.H{
					"code":    validate.Get().Code,
					"message": validate.Get().Message,
					"data":    nil,
				})
			} else if err, ok := r.(error); ok {
				cxt.JSON(http.StatusOK, gin.H{
					"code":    0,
					"message": exception.ErrorToString(err),
					"data":    nil,
				})
			} else {
				cxt.JSON(http.StatusOK, gin.H{
					"code":    0,
					"message": r,
					"data":    nil,
				})
			}
		}
	})
}

type RecoverFunc func(cxt *gin.Context, r interface{})

type Recover struct{}

func NewRecover() *Recover {
	return &Recover{}
}

func (this *Recover) Exception(hg *Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				RecoverHandle(cxt, r)
				cxt.Abort()
			}
		}()
		cxt.Next()
	}
}
