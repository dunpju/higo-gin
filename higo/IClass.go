package higo

import "reflect"

type IClass interface {
	Reflection() (reflect.Type, reflect.Value)
}
