package higo

type IClass interface {
	New() IClass
}

type Property func(class IClass)
type Propertys []Property

func (this Propertys) Apply(class IClass) {
	for _, property := range this {
		property(class)
	}
}
