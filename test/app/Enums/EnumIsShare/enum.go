package EnumIsShare

import (
	"fmt"
)

var (
	enums map[IsShare]*enum
)

const (
	Yes IsShare = 1 //是
	No IsShare = 2 //否
)

func init() {
	enums = make(map[IsShare]*enum)
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

func Enums() map[IsShare]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := IsShare(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// IsShare 是否分享
type IsShare int

func (this IsShare) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this IsShare) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this IsShare) Code() int {
	return this.get().code
}

func (this IsShare) Message() string {
	return this.get().message
}
