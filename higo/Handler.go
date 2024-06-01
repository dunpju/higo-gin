package higo

import (
	"fmt"
	"github.com/dunpju/higo-utils/utils"
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
	gid, err := utils.Runtime.GoroutineID()
	if err != nil {
		panic(err)
	}
	v, ok := this.value.Load(gid)
	if ok {
		return v.(*gin.Context)
	}
	panic(fmt.Errorf("goroutine id %d gin context empty, Cannot penetrate goroutine get gin context", gid))
}

func (this *requester) Set(ctx *gin.Context) {
	gid, err := utils.Runtime.GoroutineID()
	if err != nil {
		panic(err)
	}
	this.value.Store(gid, ctx)
}

func (this requester) Remove() {
	gid, err := utils.Runtime.GoroutineID()
	if err != nil {
		panic(err)
	}
	this.value.Delete(gid)
}

func handleConvert(handler interface{}) interface{} {
	if handle, ok := handler.(func(*gin.Context)); ok {
		return func(ctx *gin.Context) {
			defer Request.Remove()
			Request.Set(ctx)
			handle(ctx)
		}
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
