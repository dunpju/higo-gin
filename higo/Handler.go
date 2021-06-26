package higo

import (
	"github.com/dengpju/higo-utils/utils"
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
		Request = make(Req)
	})
}

type Req map[uint64]*gin.Context

func (this Req) Context() *gin.Context {
	goid := utils.GoroutineID()
	return this[goid]
}

func (this Req) Set(ctx *gin.Context) {
	goid := utils.GoroutineID()
	this[goid] = ctx
}

func (this Req) Remove() {
	goid := utils.GoroutineID()
	delete(this, goid)
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
			ctx.String(200, handle())
		}
	} else if handle, ok := handler.(func() interface{}); ok {
		return func(ctx *gin.Context) {
			defer Request.Remove()
			Request.Set(ctx)
			result := handle()
			if res, ok := result.(string); ok {
				ctx.String(200, res)
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		}
	}
	return nil
}
