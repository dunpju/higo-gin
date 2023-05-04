package EnumPreAuditState

import "github.com/dunpju/higo-enum/enum"

var e PreAuditState

func Inspect(value int) error {
	return e.Inspect(value)
}

//预审审核状态
type PreAuditState int

func (this PreAuditState) Name() string {
	return "PreAuditState"
}

func (this PreAuditState) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this PreAuditState) Message() string {
	return enum.String(this)
}

const (
	Unknown PreAuditState = 0 //未发起审核
	Waiting PreAuditState = 1 //预审待审
	Pass    PreAuditState = 2 //预审通过
	Refuse  PreAuditState = 3 //预审拒绝
)

func (this PreAuditState) Register() enum.Message {
	return make(enum.Message).
		Put(Unknown, "未发起审核").
		Put(Waiting, "预审待审").
		Put(Pass, "预审通过").
		Put(Refuse, "预审拒绝")
}
