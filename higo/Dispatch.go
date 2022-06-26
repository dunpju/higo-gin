package higo

import (
	"reflect"
)

type Dispatch struct {
	Class  IClass
	Method string
	method reflect.Value
}

func newDispatch(class IClass, method string) *Dispatch {
	return &Dispatch{Class: class.New(), Method: method, method: reflect.ValueOf(class).MethodByName(method)}
}

func (this *Dispatch) Convert(handler interface{}) interface{} {
	if handle := handleConvert(handler); handle != nil {
		return handle
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
