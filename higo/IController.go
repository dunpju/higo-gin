package higo

// 控制器接口(只要实现该接口都认为是控制器)
type IController interface {
	Controller(hg *Higo) interface{} // 构造函数
}
