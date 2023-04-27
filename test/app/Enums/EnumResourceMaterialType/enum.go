package EnumResourceMaterialType

import "github.com/dengpju/higo-enum/enum"

var e ResourceMaterialType

func Inspect(value int) error {
	return e.Inspect(value)
}

//资源材料类型
type ResourceMaterialType int

func (this ResourceMaterialType) Name() string {
	return "ResourceMaterialType"
}

func (this ResourceMaterialType) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this ResourceMaterialType) Message() string {
	return enum.String(this)
}

const (
	Word ResourceMaterialType = 1 //word
	Pdf ResourceMaterialType = 2 //pdf
)

func (this ResourceMaterialType) Register() enum.Message {
	return make(enum.Message).
	    Put(Word, "word").
	    Put(Pdf, "pdf")
}