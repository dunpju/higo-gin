package higo

import "fmt"


// 控制器接口(只要实现该接口都认为是控制器)
type IController interface {
	Controller(hg *Higo)
}

type HgController struct {
	Higo *Higo
}

func (this *HgController) Controller(hg *Higo) {
	fmt.Println("Controller")
	this.Higo = hg
}