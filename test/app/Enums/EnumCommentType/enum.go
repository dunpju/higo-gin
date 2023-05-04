package EnumCommentType

import "github.com/dunpju/higo-enum/enum"

var e CommentType

func Inspect(value int) error {
	return e.Inspect(value)
}

//评论类型
type CommentType int

func (this CommentType) Name() string {
	return "CommentType"
}

func (this CommentType) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this CommentType) Message() string {
	return enum.String(this)
}

const (
	Child            CommentType = 1 //儿童评价(tl_assessment)
	CalendarTemplate CommentType = 2 //一日模板(tl_calendar_template)
)

func (this CommentType) Register() enum.Message {
	return make(enum.Message).
		Put(Child, "儿童评价(tl_assessment)").
		Put(CalendarTemplate, "一日模板(tl_calendar_template)")
}
