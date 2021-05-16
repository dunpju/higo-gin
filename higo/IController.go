package higo

type IController interface {
	New() IClass
	Route(hg *Higo)
}
