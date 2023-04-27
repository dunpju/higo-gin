package EnumCollectType

import "github.com/dengpju/higo-enum/enum"

var e CollectType

func Inspect(value int) error {
	return e.Inspect(value)
}

//收藏类型
type CollectType int

func (this CollectType) Name() string {
	return "CollectType"
}

func (this CollectType) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this CollectType) Message() string {
	return enum.String(this)
}

const (
	Comment CollectType = 1 //评价
)

func (this CollectType) Register() enum.Message {
	return make(enum.Message).
	    Put(Comment, "评价")
}