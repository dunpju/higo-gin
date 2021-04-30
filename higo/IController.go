package higo

type IController interface {
	Self(hg *Higo) IClass
	Route(hg *Higo) *Higo
}
