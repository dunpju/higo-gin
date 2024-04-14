package EnumApplyScope

import (
	"fmt"
)

var (
	enums map[ApplyScope]*enum
)

const (
	SelfClass ApplyScope = 1 //本班级
	SelfGrade ApplyScope = 2 //本年级
	SelfRegin ApplyScope = 3 //本园所
)

func init() {
	enums = make(map[ApplyScope]*enum)
	enums[SelfClass] = newEnum(int(SelfClass), "本班级")
	enums[SelfGrade] = newEnum(int(SelfGrade), "本年级")
	enums[SelfRegin] = newEnum(int(SelfRegin), "本园所")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[ApplyScope]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := ApplyScope(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// ApplyScope 应用范围
type ApplyScope int

func (this ApplyScope) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this ApplyScope) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this ApplyScope) Code() int {
	return this.get().code
}

func (this ApplyScope) Message() string {
	return this.get().message
}
