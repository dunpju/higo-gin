package EnumSemester

import "github.com/dunpju/higo-enum/enum"

var e Semester

func Inspect(value int) error {
	return e.Inspect(value)
}

//学期
type Semester int

func (this Semester) Name() string {
	return "Semester"
}

func (this Semester) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this Semester) Message() string {
	return enum.String(this)
}

const (
	Up   Semester = 1 //上学期
	Down Semester = 2 //下学期
)

func (this Semester) Register() enum.Message {
	return make(enum.Message).
		Put(Up, "上学期").
		Put(Down, "下学期")
}
