package EnumStudentState

import "github.com/dunpju/higo-enum/enum"

var e StudentState

func Inspect(value int) error {
	return e.Inspect(value)
}

//状态
type StudentState int

func (this StudentState) Name() string {
	return "StudentState"
}

func (this StudentState) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this StudentState) Message() string {
	return enum.String(this)
}

const (
	InSchool       StudentState = 0 //在读
	FinishSchool   StudentState = 1 //毕业
	TransferSchool StudentState = 2 //转学
	StopSchool     StudentState = 3 //休学
	QuitSchool     StudentState = 4 //退学
	ExpelSchool    StudentState = 5 //开除
	Abroad         StudentState = 6 //出国
	Other          StudentState = 7 //其他
)

func (this StudentState) Register() enum.Message {
	return make(enum.Message).
		Put(InSchool, "在读").
		Put(FinishSchool, "毕业").
		Put(TransferSchool, "转学").
		Put(StopSchool, "休学").
		Put(QuitSchool, "退学").
		Put(ExpelSchool, "开除").
		Put(Abroad, "出国").
		Put(Other, "其他")
}
