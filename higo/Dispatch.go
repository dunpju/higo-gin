package higo

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type Dispatch struct {
	Class  IClass
	Method string
	method reflect.Value
}

func NewDispatch(class IClass, method string) *Dispatch {
	return &Dispatch{Class: class.New(), Method: method, method: reflect.ValueOf(class).MethodByName(method)}
}

func (this *Dispatch) Call(handler interface{}) interface{} {
	//if handle, ok := handler.(func(*gin.Context)); ok {
	//	return handle
	//}
	if handle, ok := handler.(func(*gin.Context)); ok {
		return handle
	} else if handle, ok := handler.(func() string); ok {
		return func(ctx *gin.Context) {
			ctx.String(200, handle())
		}
	}
	hRef := reflect.ValueOf(handler)
	for _, responder := range getResponderList() {
		rRef := reflect.TypeOf(responder)
		if hRef.Type().ConvertibleTo(rRef) {
			return hRef.Convert(rRef).Interface().(IResponder).Handle(this.method)
		}
	}
	panic("unknown dispatch")
}
