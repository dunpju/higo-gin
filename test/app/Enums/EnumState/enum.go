package EnumState

import "github.com/dunpju/higo-enum/enum"

var e State

func Inspect(value int) error {
	return e.Inspect(value)
}

//状态
type State int

func (this State) Name() string {
	return "State"
}

func (this State) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this State) Message() string {
	return enum.String(this)
}

const (
	Issue State = 1 //发布
	Draft State = 2 //草稿
)

func (this State) Register() enum.Message {
	return make(enum.Message).
		Put(Issue, "发布").
		Put(Draft, "草稿")
}
