package higo

import "reflect"

type Annotation interface {
	SetTag(tag reflect.StructTag)
}

var AnnotationList []Annotation

func IsAnnotation(t reflect.Type) bool {
	for _,item := range AnnotationList{
		if reflect.TypeOf(item) == t {
			return true
		}
	}
	return false
}

func init() {
	AnnotationList = make([]Annotation, 0)
	AnnotationList = append(AnnotationList,new(Value))
}