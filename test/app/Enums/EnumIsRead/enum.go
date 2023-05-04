package EnumIsRead

import "github.com/dunpju/higo-enum/enum"

var e IsRead

func Inspect(value int) error {
	return e.Inspect(value)
}

//是否阅读
type IsRead int

func (this IsRead) Name() string {
	return "IsRead"
}

func (this IsRead) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this IsRead) Message() string {
	return enum.String(this)
}

const (
	Yes IsRead = 1 //是
	No  IsRead = 2 //否
)

func (this IsRead) Register() enum.Message {
	return make(enum.Message).
		Put(Yes, "是").
		Put(No, "否")
}
