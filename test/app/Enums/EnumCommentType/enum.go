package EnumCommentType

import (
	"fmt"
)

var (
	enums map[CommentType]*enum
)

const (
	Child CommentType = 1 //儿童评价(tl_assessment)
	CalendarTemplate CommentType = 2 //一日模板(tl_calendar_template)
)

func init() {
	enums = make(map[CommentType]*enum)
	enums[Child] = newEnum(int(Child), "儿童评价(tl_assessment)")
	enums[CalendarTemplate] = newEnum(int(CalendarTemplate), "一日模板(tl_calendar_template)")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[CommentType]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := CommentType(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// CommentType 评论类型
type CommentType int

func (this CommentType) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this CommentType) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this CommentType) Code() int {
	return this.get().code
}

func (this CommentType) Message() string {
	return this.get().message
}
