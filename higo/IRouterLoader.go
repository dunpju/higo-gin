package higo

// 路由装载器(实现该接口都认为是路由)
type IRouterLoader interface {
	SetServe(serve *Serve)
	GetServe() *Serve
	Loader(hg *Higo) *Higo
}

type HiServe struct {
	*Serve
}

func NewHiServe() *HiServe {
	return &HiServe{Serve: newServe()}
}

func (this *HiServe) SetServe(serve *Serve) {
	this.Serve = serve
}

func (this *HiServe) GetServe() *Serve {
	return this.Serve
}

//func (this *HiServe) Loader(hg *Higo) *Higo {
//	return hg
//}
