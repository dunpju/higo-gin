package EnumIsShare

import "github.com/dunpju/higo-enum/enum"

var e IsShare

func Inspect(value int) error {
	return e.Inspect(value)
}

//是否分享
type IsShare int

func (this IsShare) Name() string {
	return "IsShare"
}

func (this IsShare) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this IsShare) Message() string {
	return enum.String(this)
}

const (
	Yes IsShare = 1 //是
	No  IsShare = 2 //否
)

func (this IsShare) Register() enum.Message {
	return make(enum.Message).
		Put(Yes, "是").
		Put(No, "否")
}
