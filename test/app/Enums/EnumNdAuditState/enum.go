package EnumNdAuditState

import (
	"fmt"
)

var (
	enums map[NdAuditState]*enum
)

const (
	Unknown NdAuditState = 0 //未发起审核
	Waiting NdAuditState = 1 //终审待审(预审通过自动进入此状态)
	Pass NdAuditState = 2 //终审通过
	Refuse NdAuditState = 3 //终审拒绝
)

func init() {
	enums = make(map[NdAuditState]*enum)
	enums[Unknown] = newEnum(int(Unknown), "未发起审核")
	enums[Waiting] = newEnum(int(Waiting), "终审待审(预审通过自动进入此状态)")
	enums[Pass] = newEnum(int(Pass), "终审通过")
	enums[Refuse] = newEnum(int(Refuse), "终审拒绝")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[NdAuditState]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := NdAuditState(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// NdAuditState 预审审核状态
type NdAuditState int

func (this NdAuditState) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this NdAuditState) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this NdAuditState) Code() int {
	return this.get().code
}

func (this NdAuditState) Message() string {
	return this.get().message
}
