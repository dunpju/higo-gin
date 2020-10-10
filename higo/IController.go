package higo

import "fmt"

// 控制器接口(只要实现该接口都认为是控制器)
type IController interface {
	Controller(hg *Higo) interface{}
}

type HgController struct {

}

func (this *HgController) Controller(hg *Higo) interface{} {
	return this
}

func NewController(hg *Higo, controller IController)  {
	fmt.Printf("%T\n",controller)
	controller.Controller(hg)
}