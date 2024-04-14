package EnumSectionType

import (
	"fmt"
)

var (
	enums map[SectionType]*enum
)

const (
	MorningMeeting SectionType = 1 //晨会
	SmallGroup SectionType = 2 //小组
	LargeGroup SectionType = 3 //大组
	Outdoors SectionType = 4 //户外
	Plan SectionType = 5 //计划
	Review SectionType = 6 //回顾
	CleanUp SectionType = 7 //清理
	Transition SectionType = 8 //过渡
)

func init() {
	enums = make(map[SectionType]*enum)
	enums[MorningMeeting] = newEnum(int(MorningMeeting), "晨会")
	enums[SmallGroup] = newEnum(int(SmallGroup), "小组")
	enums[LargeGroup] = newEnum(int(LargeGroup), "大组")
	enums[Outdoors] = newEnum(int(Outdoors), "户外")
	enums[Plan] = newEnum(int(Plan), "计划")
	enums[Review] = newEnum(int(Review), "回顾")
	enums[CleanUp] = newEnum(int(CleanUp), "清理")
	enums[Transition] = newEnum(int(Transition), "过渡")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[SectionType]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := SectionType(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// SectionType 环节类型
type SectionType int

func (this SectionType) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this SectionType) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this SectionType) Code() int {
	return this.get().code
}

func (this SectionType) Message() string {
	return this.get().message
}
