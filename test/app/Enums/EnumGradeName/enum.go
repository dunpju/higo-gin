package EnumGradeName

import "github.com/dengpju/higo-enum/enum"

var e GradeName

func Inspect(value string) error {
	return e.Inspect(value)
}

//幼儿园年级
type GradeName string

func (this GradeName) Name() string {
	return "GradeName"
}

func (this GradeName) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this GradeName) Message() string {
	return enum.String(this)
}

const (
	KD GradeName = "KD" //婴班
	KC GradeName = "KC" //小班
	KB GradeName = "KB" //中班
	KA GradeName = "KA" //大班
)

func (this GradeName) Register() enum.Message {
	return make(enum.Message).
	    Put(KD, "婴班").
	    Put(KC, "小班").
	    Put(KB, "中班").
	    Put(KA, "大班")
}