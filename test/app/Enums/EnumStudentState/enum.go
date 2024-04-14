package EnumStudentState

import (
	"fmt"
)

var (
	enums map[StudentState]*enum
)

const (
	InSchool StudentState = 0 //在读
	FinishSchool StudentState = 1 //毕业
	TransferSchool StudentState = 2 //转学
	StopSchool StudentState = 3 //休学
	QuitSchool StudentState = 4 //退学
	ExpelSchool StudentState = 5 //开除
	Abroad StudentState = 6 //出国
	Other StudentState = 7 //其他
)

func init() {
	enums = make(map[StudentState]*enum)
	enums[InSchool] = newEnum(int(InSchool), "在读")
	enums[FinishSchool] = newEnum(int(FinishSchool), "毕业")
	enums[TransferSchool] = newEnum(int(TransferSchool), "转学")
	enums[StopSchool] = newEnum(int(StopSchool), "休学")
	enums[QuitSchool] = newEnum(int(QuitSchool), "退学")
	enums[ExpelSchool] = newEnum(int(ExpelSchool), "开除")
	enums[Abroad] = newEnum(int(Abroad), "出国")
	enums[Other] = newEnum(int(Other), "其他")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[StudentState]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := StudentState(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// StudentState 状态
type StudentState int

func (this StudentState) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this StudentState) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this StudentState) Code() int {
	return this.get().code
}

func (this StudentState) Message() string {
	return this.get().message
}
