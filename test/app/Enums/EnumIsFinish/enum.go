package EnumIsFinish

import (
	"fmt"
)

var (
	enums map[IsFinish]*enum
)

const (
	Yes IsFinish = 1 //是
	No IsFinish = 2 //否
)

func init() {
	enums = make(map[IsFinish]*enum)
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

func Enums() map[IsFinish]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := IsFinish(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// IsFinish 是否完成
type IsFinish int

func (this IsFinish) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this IsFinish) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this IsFinish) Code() int {
	return this.get().code
}

func (this IsFinish) Message() string {
	return this.get().message
}
