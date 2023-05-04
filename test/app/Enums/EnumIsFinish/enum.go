package EnumIsFinish

import "github.com/dunpju/higo-enum/enum"

var e IsFinish

func Inspect(value int) error {
	return e.Inspect(value)
}

//是否完成
type IsFinish int

func (this IsFinish) Name() string {
	return "IsFinish"
}

func (this IsFinish) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this IsFinish) Message() string {
	return enum.String(this)
}

const (
	Yes IsFinish = 1 //是
	No  IsFinish = 2 //否
)

func (this IsFinish) Register() enum.Message {
	return make(enum.Message).
		Put(Yes, "是").
		Put(No, "否")
}
