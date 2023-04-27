package EnumDucationClassState

import "github.com/dengpju/higo-enum/enum"

var e DucationClassState

func Inspect(value int) error {
	return e.Inspect(value)
}

//状态
type DucationClassState int

func (this DucationClassState) Name() string {
	return "DucationClassState"
}

func (this DucationClassState) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this DucationClassState) Message() string {
	return enum.String(this)
}

const (
	Unknown DucationClassState = 0 //未开始
	Ing DucationClassState = 1 //进行中
	Archive DucationClassState = 2 //已归档
)

func (this DucationClassState) Register() enum.Message {
	return make(enum.Message).
	    Put(Unknown, "未开始").
	    Put(Ing, "进行中").
	    Put(Archive, "已归档")
}