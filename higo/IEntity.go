package higo

type Flag int

type IEntity interface {
	IsEdit() bool
	SetIsEdit(isEdit bool)
	SetFlag(flag Flag)
	Flag() Flag
	PriEmpty() bool
}
