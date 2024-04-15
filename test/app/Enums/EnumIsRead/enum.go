package EnumIsRead

import (
	"fmt"
)

var (
	enums map[IsRead]*enum
)

const (
	Yes IsRead = 1 //是
	No IsRead = 2 //否
)

func init() {
	enums = make(map[IsRead]*enum)
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

func Enums() map[IsRead]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := IsRead(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// IsRead 是否阅读
type IsRead int

func (this IsRead) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this IsRead) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this IsRead) Code() int {
	return this.get().code
}

func (this IsRead) Message() string {
	return this.get().message
}
