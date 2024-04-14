package EnumVerifyState

import (
	"fmt"
)

var (
	enums map[VerifyState]*enum
)

const (
	Waiting VerifyState = 1 //待审核
	Pass VerifyState = 2 //通过
	Refuse VerifyState = 3 //拒绝
)

func init() {
	enums = make(map[VerifyState]*enum)
	enums[Waiting] = newEnum(int(Waiting), "待审核")
	enums[Pass] = newEnum(int(Pass), "通过")
	enums[Refuse] = newEnum(int(Refuse), "拒绝")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[VerifyState]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := VerifyState(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// VerifyState 审核状态
type VerifyState int

func (this VerifyState) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this VerifyState) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this VerifyState) Code() int {
	return this.get().code
}

func (this VerifyState) Message() string {
	return this.get().message
}
