package event

import "reflect"

type EventHandler struct {
	fn reflect.Value
}

//执行目标函数
func (this *EventHandler) Call(params ...interface{}) interface{} {
	values := this.fn.Call(this.parseParams(params)) //调用GO 反射方法执行
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return values[0].Interface()
	}
	ret := make([]interface{}, len(values))
	for i, v := range values {
		ret[i] = v.Interface()
	}
	return ret
}

// 处理参数类型和值
func (this *EventHandler) parseParams(params []interface{}) []reflect.Value {
	parsedParams := make([]reflect.Value, len(params))
	for i, p := range params {
		if p == nil {
			//关键句
			parsedParams[i] = reflect.New(this.fn.Type().In(i)).Elem()
		} else {
			parsedParams[i] = reflect.ValueOf(p)
		}
	}
	return parsedParams
}

func NewEventHandler(fn interface{}) *EventHandler {
	getFn := reflect.ValueOf(fn)
	if getFn.Kind() != reflect.Func {
		panic("handler kind error")
	}
	return &EventHandler{fn: getFn}
}
