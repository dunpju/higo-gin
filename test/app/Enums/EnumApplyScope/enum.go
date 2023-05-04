package EnumApplyScope

import "github.com/dunpju/higo-enum/enum"

var e ApplyScope

func Inspect(value int) error {
	return e.Inspect(value)
}

//应用范围
type ApplyScope int

func (this ApplyScope) Name() string {
	return "ApplyScope"
}

func (this ApplyScope) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this ApplyScope) Message() string {
	return enum.String(this)
}

const (
	SelfClass ApplyScope = 1 //本班级
	SelfGrade ApplyScope = 2 //本年级
	SelfRegin ApplyScope = 3 //本园所
)

func (this ApplyScope) Register() enum.Message {
	return make(enum.Message).
		Put(SelfClass, "本班级").
		Put(SelfGrade, "本年级").
		Put(SelfRegin, "本园所")
}
