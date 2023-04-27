package EnumIsRecheck

import "github.com/dengpju/higo-enum/enum"

var e IsRecheck

func Inspect(value int) error {
	return e.Inspect(value)
}

//是否复核
type IsRecheck int

func (this IsRecheck) Name() string {
	return "IsRecheck"
}

func (this IsRecheck) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this IsRecheck) Message() string {
	return enum.String(this)
}

const (
	Yes IsRecheck = 1 //是
	No IsRecheck = 2 //否
)

func (this IsRecheck) Register() enum.Message {
	return make(enum.Message).
	    Put(Yes, "是").
	    Put(No, "否")
}