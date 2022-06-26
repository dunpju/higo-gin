package higo

import (
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	Request     requester
	onceRequest sync.Once
)

const goroutineID = "GoroutineID"

func init() {
	onceRequest.Do(func() {
		Request = requester{value: &sync.Map{}}
	})
}

type requester struct {
	value *sync.Map
}

func (this *requester) Context() *gin.Context {
	goid, err := utils.Runtime.GoroutineID()
	if err != nil {
		panic(err)
	}
	v, ok := this.value.Load(goid)
	if ok {
		return v.(*gin.Context)
	}
	panic(fmt.Errorf("goroutine id %d gin context empty, Cannot penetrate goroutine get gin context", goid))
}

func (this *requester) Set(ctx *gin.Context) {
	goid, err := utils.Runtime.GoroutineID()
	if err != nil {
		panic(err)
	}
	this.value.Store(goid, ctx)
}

func (this requester) Remove() {
	goid, err := utils.Runtime.GoroutineID()
	if err != nil {
		panic(err)
	}
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
