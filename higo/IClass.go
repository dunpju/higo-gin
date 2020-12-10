package higo

import "reflect"

// 类接口(只要实现该接口都认为是类)
type IClass interface {
	Class() (reflect.Type, reflect.Value)
}
