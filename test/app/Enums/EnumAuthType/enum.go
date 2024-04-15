package EnumAuthType

import (
	"fmt"
)

var (
	enums map[AuthType]*enum
)

const (
	Action AuthType = 1 //功能
	Data AuthType = 2 //数据
)

func init() {
	enums = make(map[AuthType]*enum)
	enums[Action] = newEnum(int(Action), "功能")
	enums[Data] = newEnum(int(Data), "数据")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[AuthType]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := AuthType(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// AuthType 权限类型
type AuthType int

func (this AuthType) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this AuthType) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this AuthType) Code() int {
	return this.get().code
}

func (this AuthType) Message() string {
	return this.get().message
}
