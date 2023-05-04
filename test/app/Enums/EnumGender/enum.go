package EnumGender

import "github.com/dunpju/higo-enum/enum"

var e Gender

func Inspect(value int) error {
	return e.Inspect(value)
}

//是否正在使用
type Gender int

func (this Gender) Name() string {
	return "Gender"
}

func (this Gender) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this Gender) Message() string {
	return enum.String(this)
}

const (
	Unknown Gender = 0 //未知
	Male    Gender = 1 //男
	Female  Gender = 2 //女
)

func (this Gender) Register() enum.Message {
	return make(enum.Message).
		Put(Unknown, "未知").
		Put(Male, "男").
		Put(Female, "女")
}
