package higo

import "fmt"

//生成器接口
type IBuilder interface {
	Construct()
}

type Director struct {
	builder IBuilder
}

func NewDirector(builders ...IBuilder) {
	for _,builder := range builders{
		builder.Construct()
	}
}

func (d *Director) Construct()  {
	fmt.Printf("%p\n", d)
	fmt.Printf("%p\n", d.builder)
}

