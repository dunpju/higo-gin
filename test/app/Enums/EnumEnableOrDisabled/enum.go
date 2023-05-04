package EnumEnableOrDisabled

import "github.com/dunpju/higo-enum/enum"

var e EnableOrDisabled

func Inspect(value int) error {
	return e.Inspect(value)
}

//启用禁用
type EnableOrDisabled int

func (this EnableOrDisabled) Name() string {
	return "EnableOrDisabled"
}

func (this EnableOrDisabled) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this EnableOrDisabled) Message() string {
	return enum.String(this)
}

const (
	Enable   EnableOrDisabled = 1 //启用
	Disabled EnableOrDisabled = 2 //禁用
)

func (this EnableOrDisabled) Register() enum.Message {
	return make(enum.Message).
		Put(Enable, "启用").
		Put(Disabled, "禁用")
}
