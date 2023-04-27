package EnumIsAbnormal

import "github.com/dengpju/higo-enum/enum"

var e IsAbnormal

func Inspect(value int) error {
	return e.Inspect(value)
}

//分数异常
type IsAbnormal int

func (this IsAbnormal) Name() string {
	return "IsAbnormal"
}

func (this IsAbnormal) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this IsAbnormal) Message() string {
	return enum.String(this)
}

const (
	Normal IsAbnormal = 0 //正常
	Inc IsAbnormal = 1 //上升
	Dec IsAbnormal = 2 //下降
)

func (this IsAbnormal) Register() enum.Message {
	return make(enum.Message).
	    Put(Normal, "正常").
	    Put(Inc, "上升").
	    Put(Dec, "下降")
}