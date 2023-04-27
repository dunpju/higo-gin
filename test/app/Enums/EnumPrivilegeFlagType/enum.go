package EnumPrivilegeFlagType

import "github.com/dengpju/higo-enum/enum"

var e PrivilegeFlagType

func Inspect(value int) error {
	return e.Inspect(value)
}

//权限标签类型
type PrivilegeFlagType int

func (this PrivilegeFlagType) Name() string {
	return "PrivilegeFlagType"
}

func (this PrivilegeFlagType) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this PrivilegeFlagType) Message() string {
	return enum.String(this)
}

const (
	Action PrivilegeFlagType = 1 //功能
	Menu PrivilegeFlagType = 2 //菜单
)

func (this PrivilegeFlagType) Register() enum.Message {
	return make(enum.Message).
	    Put(Action, "功能").
	    Put(Menu, "菜单")
}