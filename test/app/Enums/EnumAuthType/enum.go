package EnumAuthType

import "github.com/dunpju/higo-enum/enum"

var e AuthType

func Inspect(value int) error {
	return e.Inspect(value)
}

//权限类型
type AuthType int

func (this AuthType) Name() string {
	return "AuthType"
}

func (this AuthType) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this AuthType) Message() string {
	return enum.String(this)
}

const (
	Action AuthType = 1 //功能
	Data   AuthType = 2 //数据
)

func (this AuthType) Register() enum.Message {
	return make(enum.Message).
		Put(Action, "功能").
		Put(Data, "数据")
}
