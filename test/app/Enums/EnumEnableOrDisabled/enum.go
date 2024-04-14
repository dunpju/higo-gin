package EnumEnableOrDisabled

import (
	"fmt"
)

var (
	enums map[EnableOrDisabled]*enum
)

const (
	Enable EnableOrDisabled = 1 //启用
	Disabled EnableOrDisabled = 2 //禁用
)

func init() {
	enums = make(map[EnableOrDisabled]*enum)
	enums[Enable] = newEnum(int(Enable), "启用")
	enums[Disabled] = newEnum(int(Disabled), "禁用")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[EnableOrDisabled]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := EnableOrDisabled(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// EnableOrDisabled 启用禁用
type EnableOrDisabled int

func (this EnableOrDisabled) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this EnableOrDisabled) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this EnableOrDisabled) Code() int {
	return this.get().code
}

func (this EnableOrDisabled) Message() string {
	return this.get().message
}
