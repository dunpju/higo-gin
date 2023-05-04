package EnumVersionIsUsed

import "github.com/dunpju/higo-enum/enum"

var e VersionIsUsed

func Inspect(value int) error {
	return e.Inspect(value)
}

//是否正在使用
type VersionIsUsed int

func (this VersionIsUsed) Name() string {
	return "VersionIsUsed"
}

func (this VersionIsUsed) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this VersionIsUsed) Message() string {
	return enum.String(this)
}

const (
	Yes VersionIsUsed = 1 //是
	No  VersionIsUsed = 2 //否
)

func (this VersionIsUsed) Register() enum.Message {
	return make(enum.Message).
		Put(Yes, "是").
		Put(No, "否")
}
