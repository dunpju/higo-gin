package EnumIsAbnormal

import (
	"fmt"
)

var (
	enums map[IsAbnormal]*enum
)

const (
	Normal IsAbnormal = 0 //正常
	Inc IsAbnormal = 1 //上升
	Dec IsAbnormal = 2 //下降
)

func init() {
	enums = make(map[IsAbnormal]*enum)
	enums[Normal] = newEnum(int(Normal), "正常")
	enums[Inc] = newEnum(int(Inc), "上升")
	enums[Dec] = newEnum(int(Dec), "下降")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[IsAbnormal]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := IsAbnormal(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// IsAbnormal 分数异常
type IsAbnormal int

func (this IsAbnormal) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this IsAbnormal) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this IsAbnormal) Code() int {
	return this.get().code
}

func (this IsAbnormal) Message() string {
	return this.get().message
}
