package EnumIsSuper

import (
	"fmt"
)

var (
	enums map[IsSuper]*enum
)

const (
	Yes IsSuper = 1 //是
	No IsSuper = 2 //否
)

func init() {
	enums = make(map[IsSuper]*enum)
	enums[Yes] = newEnum(int(Yes), "是")
	enums[No] = newEnum(int(No), "否")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[IsSuper]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := IsSuper(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// IsSuper 是否超级管理员
type IsSuper int

func (this IsSuper) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this IsSuper) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this IsSuper) Code() int {
	return this.get().code
}

func (this IsSuper) Message() string {
	return this.get().message
}
