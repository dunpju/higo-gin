package EnumDeleteAuditState

import (
	"fmt"
)

var (
	enums map[DeleteAuditState]*enum
)

const (
	NotDelete DeleteAuditState = 0 //未发起删除
	Waiting DeleteAuditState = 1 //待审核
	Pass DeleteAuditState = 2 //审核通过
	NotPass DeleteAuditState = 3 //审核未通过
)

func init() {
	enums = make(map[DeleteAuditState]*enum)
	enums[NotDelete] = newEnum(int(NotDelete), "未发起删除")
	enums[Waiting] = newEnum(int(Waiting), "待审核")
	enums[Pass] = newEnum(int(Pass), "审核通过")
	enums[NotPass] = newEnum(int(NotPass), "审核未通过")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[DeleteAuditState]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := DeleteAuditState(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// DeleteAuditState 删除审核状态
type DeleteAuditState int

func (this DeleteAuditState) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this DeleteAuditState) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this DeleteAuditState) Code() int {
	return this.get().code
}

func (this DeleteAuditState) Message() string {
	return this.get().message
}
