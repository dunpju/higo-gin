package higo

import (
	"fmt"
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-logger/logger"
	"github.com/dengpju/higo-throw/throw"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
				logger.Logrus.Info(fmt.Sprintf("Recover Value %v", r))
				logger.Logrus.Info(fmt.Sprintf("Recover Type %T", r))
				//打印错误堆栈信息
				//debug.PrintStack()
				// 输出换行debug调用栈
				logger.PrintlnStack()
				//封装通用json返回
				if h, ok := r.(gin.H); ok {
					c.JSON(http.StatusOK, h)
				} else {
					if msg, ok := r.(*code.Code); ok {
						c.JSON(http.StatusOK, gin.H{
							"code": msg.Code,
							"msg":  msg.Message,
							"data": nil,
						})
					} else {
						c.JSON(http.StatusOK, gin.H{
							"code": 0,
							"msg":  throw.ErrorToString(r),
							"data": nil,
						})
					}
				}
				//终止
				c.Abort()
			}
		}()
		//继续
		c.Next()
	}
}
