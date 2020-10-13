package higo

import "fmt"

//生成器接口
type IBuilder interface {
	Construct()
}

type Director struct {
	builder IBuilder
}

func NewDirector(builder IBuilder) *Director {
	return &Director{
		builder: builder,
	}
}

func (d *Director) Construct()  {
	fmt.Printf("%p\n", d)
	fmt.Printf("%p\n", d.builder)
}

