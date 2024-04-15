package EnumNoticeType

import (
	"fmt"
)

var (
	enums map[NoticeType]*enum
)

const (
	System NoticeType = 1 //系统通知
	Comment NoticeType = 2 //评论
	Fans NoticeType = 3 //粉丝
	Praise NoticeType = 4 //点赞
	AssessmentCorvidae NoticeType = 5 //评价待完善
	ShareRemind NoticeType = 6 //分享提醒
	Section NoticeType = 7 //环节设置
)

func init() {
	enums = make(map[NoticeType]*enum)
	enums[System] = newEnum(int(System), "系统通知")
	enums[Comment] = newEnum(int(Comment), "评论")
	enums[Fans] = newEnum(int(Fans), "粉丝")
	enums[Praise] = newEnum(int(Praise), "点赞")
	enums[AssessmentCorvidae] = newEnum(int(AssessmentCorvidae), "评价待完善")
	enums[ShareRemind] = newEnum(int(ShareRemind), "分享提醒")
	enums[Section] = newEnum(int(Section), "环节设置")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[NoticeType]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := NoticeType(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// NoticeType 消息类型
type NoticeType int

func (this NoticeType) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this NoticeType) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this NoticeType) Code() int {
	return this.get().code
}

func (this NoticeType) Message() string {
	return this.get().message
}
