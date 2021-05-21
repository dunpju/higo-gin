package higo

import (
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
	hRef := reflect.ValueOf(handler)
	for _, r := range getResponderList() {
		rRef := reflect.TypeOf(r)
		if hRef.Type().ConvertibleTo(rRef) {
			return hRef.Convert(rRef).Interface().(IResponder).Handle(this.method)
		}
	}
	panic("unknown dispatch")
}
