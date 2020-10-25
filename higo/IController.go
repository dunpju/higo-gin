package higo


// 控制器接口(只要实现该接口都认为是控制器)
type IController interface {
	Controller() *Higo
}

type HgController struct {
	Hg int
}

func NewHgController() *HgController {
	return &HgController{Hg:35}
}

func NewController(controller IController) {
	controller.Controller()
}

/**
func (this *HgController) Controller() *Higo {
	fmt.Println("Controller")
	fmt.Printf("%p\n",this)
	fmt.Println(this.Higo)
	return this.Higo
}*/