package injector

import (
	"fmt"
	"reflect"
)

type BeanFactory struct {
	beans []interface{}
}

func NewBeanFactory() *BeanFactory {
	bf := &BeanFactory{beans: make([]interface{}, 0)}
	bf.beans = append(bf.beans, bf)
	return &BeanFactory{}
}

func (this *BeanFactory) getBean(t reflect.Type) interface{} {
	fmt.Println(this.beans)
	for _, p := range this.beans {
		fmt.Println("t:",t)
		fmt.Println("pt:",reflect.TypeOf(p))
		if t == reflect.TypeOf(p) {
			fmt.Println("find")
			return p
		}
	}
	return nil
}

func (this *BeanFactory) SetBean(bean ...interface{}) {
	this.beans = append(this.beans, bean...)
}

func (this *BeanFactory) GetBean(bean interface{}) interface{} {
	return this.getBean(reflect.TypeOf(bean))
}

func (this *BeanFactory) inject(object interface{}) {
	vObject := reflect.ValueOf(object)
	if vObject.Kind() == reflect.Ptr {
		vObject = vObject.Elem()
	}
	for i := 0; i < vObject.NumField(); i++ {
		f := vObject.Field(i)
		if f.Kind() != reflect.Ptr || !f.IsNil() {
			continue
		}
		if p := this.getBean(f.Type()); p != nil && f.CanInterface() {
			f.Set(reflect.New(f.Type().Elem()))
			f.Elem().Set(reflect.ValueOf(p).Elem())
		}
	}
}

func (this *BeanFactory) Inject(bean interface{}) {
	bv := reflect.ValueOf(bean).Elem()
	bt := reflect.TypeOf(bean).Elem()
	fmt.Printf("%T\n",bean)
	for i := 0; i < bv.NumField(); i++ {
		f := bv.Field(i)
		if f.Kind() != reflect.Ptr || !f.IsNil() {
			continue
		}
		if IsAnnotation(f.Type()) {
			f.Set(reflect.New(f.Type().Elem()))
			f.Interface().(Annotation).SetTag(bt.Field(i).Tag)
			this.inject(f.Interface())
			continue
		}
		if p := this.getBean(f.Type()); p != nil {
			f.Set(reflect.New(f.Type().Elem()))
			f.Elem().Set(reflect.ValueOf(p).Elem())
			fmt.Println("f:",f)
			fmt.Println("bean:",bean)
			fmt.Println("bean ptr:",&bean)
		}
	}
}
