package EnumSectionType

import "github.com/dunpju/higo-enum/enum"

var e SectionType

func Inspect(value int) error {
	return e.Inspect(value)
}

//环节类型
type SectionType int

func (this SectionType) Name() string {
	return "SectionType"
}

func (this SectionType) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this SectionType) Message() string {
	return enum.String(this)
}

const (
	MorningMeeting SectionType = 1 //晨会
	SmallGroup     SectionType = 2 //小组
	LargeGroup     SectionType = 3 //大组
	Outdoors       SectionType = 4 //户外
	Plan           SectionType = 5 //计划
	Review         SectionType = 6 //回顾
	CleanUp        SectionType = 7 //清理
	Transition     SectionType = 8 //过渡
)

func (this SectionType) Register() enum.Message {
	return make(enum.Message).
		Put(MorningMeeting, "晨会").
		Put(SmallGroup, "小组").
		Put(LargeGroup, "大组").
		Put(Outdoors, "户外").
		Put(Plan, "计划").
		Put(Review, "回顾").
		Put(CleanUp, "清理").
		Put(Transition, "过渡")
}
