package EnumIsSuper

import "github.com/dengpju/higo-enum/enum"

var e IsSuper

func Inspect(value int) error {
	return e.Inspect(value)
}

//是否超级管理员
type IsSuper int

func (this IsSuper) Name() string {
	return "IsSuper"
}

func (this IsSuper) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this IsSuper) Message() string {
	return enum.String(this)
}

const (
	Yes IsSuper = 1 //是
	No IsSuper = 2 //否
)

func (this IsSuper) Register() enum.Message {
	return make(enum.Message).
	    Put(Yes, "是").
	    Put(No, "否")
}