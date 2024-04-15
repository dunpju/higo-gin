package EnumPreAuditState

import (
	"fmt"
)

var (
	enums map[PreAuditState]*enum
)

const (
	Unknown PreAuditState = 0 //未发起审核
	Waiting PreAuditState = 1 //预审待审
	Pass PreAuditState = 2 //预审通过
	Refuse PreAuditState = 3 //预审拒绝
)

func init() {
	enums = make(map[PreAuditState]*enum)
	enums[Unknown] = newEnum(int(Unknown), "未发起审核")
	enums[Waiting] = newEnum(int(Waiting), "预审待审")
	enums[Pass] = newEnum(int(Pass), "预审通过")
	enums[Refuse] = newEnum(int(Refuse), "预审拒绝")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[PreAuditState]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := PreAuditState(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// PreAuditState 预审审核状态
type PreAuditState int

func (this PreAuditState) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this PreAuditState) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this PreAuditState) Code() int {
	return this.get().code
}

func (this PreAuditState) Message() string {
	return this.get().message
}
