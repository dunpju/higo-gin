package higo

type IController interface {
	Route(hg *Higo) *Higo
}
