package higo

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"sync"
)

var responderList []Responder
var onceRespList sync.Once

type Responder interface {
	RespondTo() gin.HandlerFunc
}

func getResponderList() []Responder {
	onceRespList.Do(func() {
		responderList = []Responder{
			(StringResponder)(nil),
			(JsonResponder)(nil),
			(ModelResponder)(nil),
			(ModelsResponder)(nil),
		}
	})
	return responderList
}

type StringResponder func(*gin.Context) string

func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(200, getSyncHandler().handler(this, context).(string))
	}
}

type Json interface{}
type JsonResponder func(*gin.Context) Json

func (this JsonResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, getSyncHandler().handler(this, context))
	}
}

// 转换
func Convert(handler interface{}) gin.HandlerFunc {
	hRef := reflect.ValueOf(handler)
	for _, r := range getResponderList() {
		rRef := reflect.TypeOf(r)
		if hRef.Type().ConvertibleTo(rRef) {
			return hRef.Convert(rRef).Interface().(Responder).RespondTo()
		}
	}
	return nil
}

var syncHandler *SyncHandler

func getSyncHandler() *SyncHandler {
	Once.Do(func() {
		syncHandler = &SyncHandler{}
	})
	return syncHandler
}

type SyncHandler struct {
	context []IContext
}

func (this *SyncHandler) handler(responder Responder, ctx *gin.Context) interface{} {
	var ret interface{}
	if s1, ok := responder.(StringResponder); ok {
		ret = s1(ctx)
	}
	if s2, ok := responder.(JsonResponder); ok {
		ret = s2(ctx)
	}
	return ret
}

type ModelResponder func(*gin.Context) Model

func (this ModelResponder) RespondTo() gin.HandlerFunc  {
	return func(context *gin.Context) {
		context.JSON(200, this(context))
	}
}

type ModelsResponder func(*gin.Context) Models

func (this ModelsResponder) RespondTo() gin.HandlerFunc  {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Content-typ","application/json")
		_, _ = context.Writer.WriteString(string(this(context)))
	}
}