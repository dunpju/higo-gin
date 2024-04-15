package EnumPrivilegeFlagType

import (
	"fmt"
)

var (
	enums map[PrivilegeFlagType]*enum
)

const (
	Action PrivilegeFlagType = 1 //功能
	Menu PrivilegeFlagType = 2 //菜单
)

func init() {
	enums = make(map[PrivilegeFlagType]*enum)
	enums[Action] = newEnum(int(Action), "功能")
	enums[Menu] = newEnum(int(Menu), "菜单")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[PrivilegeFlagType]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := PrivilegeFlagType(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// PrivilegeFlagType 权限标签类型
type PrivilegeFlagType int

func (this PrivilegeFlagType) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this PrivilegeFlagType) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this PrivilegeFlagType) Code() int {
	return this.get().code
}

func (this PrivilegeFlagType) Message() string {
	return this.get().message
}
