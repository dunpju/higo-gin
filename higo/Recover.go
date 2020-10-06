package higo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 全局异常
type Recover struct {}

// 构造函数
func NewRecover() *Recover {
	return &Recover{}
}

func (this *Recover) RuntimeException(hg *Higo) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				Logrus.Info(fmt.Sprintf("Value %v\n", r))
				Logrus.Info(fmt.Sprintf("Type %T\n", r))
				//打印错误堆栈信息
				//debug.PrintStack()
				// 输出换行debug调用栈
				PrintlnStack()
				//封装通用json返回
				if _, ok := r.(gin.H); ok {// 断言类型
					c.JSON(http.StatusOK, r)
				} else {
					// 指针类型所以必须加指针才能取出值
					if _, ok := r.(*CodeMsg); ok {
						c.JSON(http.StatusOK, gin.H{
							"code": (r.(*CodeMsg)).Code,
							"msg":  (r.(*CodeMsg)).Msg,
							"data": nil,
						})
					} else {
						c.JSON(http.StatusOK, gin.H{
							"code": 0,
							"msg":  ErrorToString(r),
							"data": nil,
						})
					}
				}
				//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
				c.Abort()
			}
		}()
		//加载完 defer recover，继续后续接口调用
		c.Next()
	}
}

