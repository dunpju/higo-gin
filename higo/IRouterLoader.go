package higo

// 路由装载器(实现该接口都认为是路由)
type IRouterLoader interface {
	Loader(hg *Higo) *Higo
}