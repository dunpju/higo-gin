package higo

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"sync"
)

var (
	responderList      []IResponder
	onceRespList       sync.Once
	getSyncHandlerOnce sync.Once
	syncHandler        *SyncHandler
)

type IResponder interface {
	RespondTo() gin.HandlerFunc
	Handle(method reflect.Value) interface{}
}

func getResponderList() []IResponder {
	onceRespList.Do(func() {
		responderList = []IResponder{
			(StringResponder)(nil),
			(JsonResponder)(nil),
			(ModelResponder)(nil),
			(ModelsResponder)(nil),
			(WebsocketResponder)(nil),
		}
	})
	return responderList
}

// 转换
func Convert(handler interface{}) gin.HandlerFunc {
	if handle, ok := handler.(func(*gin.Context)); ok {
		return handle
	}
	hRef := reflect.ValueOf(handler)
	for _, r := range getResponderList() {
		rRef := reflect.TypeOf(r)
		if hRef.Type().ConvertibleTo(rRef) {
			return hRef.Convert(rRef).Interface().(IResponder).RespondTo()
		}
	}
	return nil
}

func getSyncHandler() *SyncHandler {
	getSyncHandlerOnce.Do(func() {
		syncHandler = &SyncHandler{}
	})
	return syncHandler
}

func methodCall(ctx *gin.Context, method reflect.Value) interface{} {
	params := make([]reflect.Value, 0)
	params = append(params, reflect.ValueOf(ctx))
	callRet := method.Call(params)
	if callRet != nil && len(callRet) == 1 {
		return callRet[0].Interface()
	}
	panic("method call error")
}

type SyncHandler struct {
	context []IContext
}

func (this *SyncHandler) handler(responder IResponder, ctx *gin.Context) interface{} {
	var ret interface{}
	if s1, ok := responder.(StringResponder); ok {
		ret = s1(ctx)
	}
	if s2, ok := responder.(JsonResponder); ok {
		ret = s2(ctx)
	}
	return ret
}

type StringResponder func(*gin.Context) string

func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(200, getSyncHandler().handler(this, ctx).(string))
	}
}

func (this StringResponder) Handle(method reflect.Value) interface{} {
	return func(ctx *gin.Context) string {
		return methodCall(ctx, method).(string)
	}
}

type Json interface{}
type JsonResponder func(*gin.Context) Json

func (this JsonResponder) RespondTo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, getSyncHandler().handler(this, ctx))
	}
}

func (this JsonResponder) Handle(method reflect.Value) interface{} {
	return func(ctx *gin.Context) Json {
		return methodCall(ctx, method).(Json)
	}
}

type ModelResponder func(*gin.Context) Model

func (this ModelResponder) RespondTo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, this(ctx))
	}
}

func (this ModelResponder) Handle(method reflect.Value) interface{} {
	return func(ctx *gin.Context) Model {
		return methodCall(ctx, method).(Model)
	}
}

type ModelsResponder func(*gin.Context) Models

func (this ModelsResponder) RespondTo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-typ", "application/json")
		_, err := ctx.Writer.WriteString(string(this(ctx)))
		if err != nil {
			panic(err)
		}
	}
}

func (this ModelsResponder) Handle(method reflect.Value) interface{} {
	return func(ctx *gin.Context) Models {
		return methodCall(ctx, method).(Models)
	}
}

type WebsocketResponder func(*gin.Context) WsWriteMessage

func (this WebsocketResponder) RespondTo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		this(ctx)
	}
}

func (this WebsocketResponder) Handle(method reflect.Value) interface{} {
	return func(ctx *gin.Context) WsWriteMessage {
		return methodCall(ctx, method).(WsWriteMessage)
	}
}
