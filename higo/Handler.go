package higo

import (
	"fmt"
	"github.com/dengpju/higo-utils/utils/runtimeutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	Request     Req
	onceRequest sync.Once
)

func init() {
	onceRequest.Do(func() {
		Request = Req{value: &sync.Map{}}
	})
}

type Req struct {
	value *sync.Map
}

func (this Req) Context() *gin.Context {
	goid := runtimeutil.GoroutineID()
	v, ok := this.value.Load(goid)
	if ok {
		return v.(*gin.Context)
	}
	panic(fmt.Errorf("goroutine id %d gin context empty", goid))
}

func (this Req) Set(ctx *gin.Context) {
	goid := runtimeutil.GoroutineID()
	this.value.Store(goid, ctx)
}

func (this Req) Remove() {
	goid := runtimeutil.GoroutineID()
	this.value.Delete(goid)
}

func handleConvert(handler interface{}) interface{} {
	if handle, ok := handler.(func(*gin.Context)); ok {
		return handle
	} else if handle, ok := handler.(func()); ok {
		return func(ctx *gin.Context) {
			defer Request.Remove()
			Request.Set(ctx)
			handle()
		}
	} else if handle, ok := handler.(func() string); ok {
		return func(ctx *gin.Context) {
			defer Request.Remove()
			Request.Set(ctx)
			ctx.String(http.StatusOK, handle())
		}
	} else if handle, ok := handler.(func() interface{}); ok {
		return func(ctx *gin.Context) {
			defer Request.Remove()
			Request.Set(ctx)
			result := handle()
			if res, ok := result.(string); ok {
				ctx.String(http.StatusOK, res)
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		}
	}
	return nil
}
