package higo

import (
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-logger/logger"
	"github.com/dengpju/higo-throw/throw"
	"github.com/dengpju/higo-utils/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"sync"
)

var (
	//Recover处理函数(可自定义)
	RecoverHandlerFunc RecoverFunc
	recoverOnce        sync.Once
)

func init() {
	recoverOnce.Do(func() {
		//初始化Recover处理函数
		RecoverHandlerFunc = func(c *gin.Context, r interface{}) {
			//打印错误堆栈信息
			debug.PrintStack()
			// 输出换行debug调用栈
			go logger.PrintlnStack()
			//封装通用json返回
			if h, ok := r.(gin.H); ok {
				c.JSON(http.StatusOK, h)
			} else if msg, ok := r.(*code.Code); ok {
				c.JSON(http.StatusOK, gin.H{
					"code":    msg.Code,
					"message": msg.Message,
					"data":    nil,
				})
			} else if MapString, ok := r.(utils.MapString); ok {
				c.JSON(http.StatusOK, MapString)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    0,
					"message": throw.ErrorToString(r),
					"data":    nil,
				})
			}
		}
	})
}

type RecoverFunc func(c *gin.Context, r interface{})

// 全局异常
type Recover struct{}

// 构造函数
func NewRecover() *Recover {
	return &Recover{}
}

func (this *Recover) Exception(hg *Higo) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				RecoverHandlerFunc(c, r) //执行Recover处理函数
				c.Abort()                //终止
			}
		}()
		c.Next() //继续
	}
}
