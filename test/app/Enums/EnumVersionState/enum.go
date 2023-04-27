package EnumVersionState

import "github.com/dengpju/higo-enum/enum"

var e VersionState

func Inspect(value int) error {
	return e.Inspect(value)
}

//版本状态
type VersionState int

func (this VersionState) Name() string {
	return "VersionState"
}

func (this VersionState) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this VersionState) Message() string {
	return enum.String(this)
}

const (
	Issue VersionState = 1 //发布
	Waiting VersionState = 2 //待发布
)

func (this VersionState) Register() enum.Message {
	return make(enum.Message).
	    Put(Issue, "发布").
	    Put(Waiting, "待发布")
}