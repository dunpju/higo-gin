package injector

import "reflect"

var BeanFactory *BeanFactoryImpl

func init()  {
	BeanFactory=NewBeanFactory()
}

type BeanFactoryImpl struct {
	beanMapper BeanMapper
}

func (this *BeanFactoryImpl)Set(value ...interface{})  {
	if value == nil || len(value)==0{
		return
	}
	for _,v := range value{
		this.beanMapper.add(v)
	}
}

func (this *BeanFactoryImpl)Get(v interface{}) interface{} {
	if v == nil{
		return nil
	}
	value := this.beanMapper.get(v)
	if value.IsValid() {
		return value.Interface()
	}
	return nil
}

// 处理依赖注入
func (this *BeanFactoryImpl)Apply(bean interface{})  {
	if bean==nil {
		return
	}
	v:=reflect.ValueOf(bean)
	if v.Kind()==reflect.Ptr {
		v=v.Elem()
	}
	if v.Kind()!= reflect.Struct {
		return
	}
	for i:=0;i<v.NumField() ;i++  {
		field:=v.Type().Field(i)
		if v.Field(i).CanSet() && field.Tag.Get("inject") != ""{
			if value:=this.Get(field.Type);value!=nil{
				v.Field(i).Set(reflect.ValueOf(value))
			}
		}
	}
}

func NewBeanFactory() *BeanFactoryImpl {
	return &BeanFactoryImpl{beanMapper: make(BeanMapper)}
}
