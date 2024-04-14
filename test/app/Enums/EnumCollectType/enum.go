package EnumCollectType

import (
	"fmt"
)

var (
	enums map[CollectType]*enum
)

const (
	Comment CollectType = 1 //评价
)

func init() {
	enums = make(map[CollectType]*enum)
	enums[Comment] = newEnum(int(Comment), "评价")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[CollectType]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := CollectType(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// CollectType 收藏类型
type CollectType int

func (this CollectType) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this CollectType) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this CollectType) Code() int {
	return this.get().code
}

func (this CollectType) Message() string {
	return this.get().message
}
