package higo

/**
// 控制器接口(只要实现该接口都认为是控制器)
type IController interface {
	Controller() interface{}
}

type HgController struct {

}

func (this *HgController) Controller() interface{} {
	fmt.Println("Controller")
	return "Controller"
}

func NewController(hg *Higo, controller IController)  {
	fmt.Printf("%T\n",controller)
	controller.Controller()
}*/