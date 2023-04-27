package EnumNdAuditState

import "github.com/dengpju/higo-enum/enum"

var e NdAuditState

func Inspect(value int) error {
	return e.Inspect(value)
}

//预审审核状态
type NdAuditState int

func (this NdAuditState) Name() string {
	return "NdAuditState"
}

func (this NdAuditState) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this NdAuditState) Message() string {
	return enum.String(this)
}

const (
	Unknown NdAuditState = 0 //未发起审核
	Waiting NdAuditState = 1 //终审待审(预审通过自动进入此状态)
	Pass NdAuditState = 2 //终审通过
	Refuse NdAuditState = 3 //终审拒绝
)

func (this NdAuditState) Register() enum.Message {
	return make(enum.Message).
	    Put(Unknown, "未发起审核").
	    Put(Waiting, "终审待审(预审通过自动进入此状态)").
	    Put(Pass, "终审通过").
	    Put(Refuse, "终审拒绝")
}