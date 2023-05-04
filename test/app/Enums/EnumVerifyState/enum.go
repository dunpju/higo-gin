package EnumVerifyState

import "github.com/dunpju/higo-enum/enum"

var e VerifyState

func Inspect(value int) error {
	return e.Inspect(value)
}

//审核状态
type VerifyState int

func (this VerifyState) Name() string {
	return "VerifyState"
}

func (this VerifyState) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this VerifyState) Message() string {
	return enum.String(this)
}

const (
	Waiting VerifyState = 1 //待审核
	Pass    VerifyState = 2 //通过
	Refuse  VerifyState = 3 //拒绝
)

func (this VerifyState) Register() enum.Message {
	return make(enum.Message).
		Put(Waiting, "待审核").
		Put(Pass, "通过").
		Put(Refuse, "拒绝")
}
