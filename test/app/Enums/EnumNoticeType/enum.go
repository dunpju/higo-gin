package EnumNoticeType

import "github.com/dunpju/higo-enum/enum"

var e NoticeType

func Inspect(value int) error {
	return e.Inspect(value)
}

//消息类型
type NoticeType int

func (this NoticeType) Name() string {
	return "NoticeType"
}

func (this NoticeType) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this NoticeType) Message() string {
	return enum.String(this)
}

const (
	System             NoticeType = 1 //系统通知
	Comment            NoticeType = 2 //评论
	Fans               NoticeType = 3 //粉丝
	Praise             NoticeType = 4 //点赞
	AssessmentCorvidae NoticeType = 5 //评价待完善
	ShareRemind        NoticeType = 6 //分享提醒
	Section            NoticeType = 7 //环节设置
)

func (this NoticeType) Register() enum.Message {
	return make(enum.Message).
		Put(System, "系统通知").
		Put(Comment, "评论").
		Put(Fans, "粉丝").
		Put(Praise, "点赞").
		Put(AssessmentCorvidae, "评价待完善").
		Put(ShareRemind, "分享提醒").
		Put(Section, "环节设置")
}
