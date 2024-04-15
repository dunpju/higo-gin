package EnumCategory

import (
	"fmt"
)

var (
	enums map[Category]*enum
)

const (
	Dept Category = 1 //部门
	Company Category = 2 //公司
	Team Category = 3 //团队
	Group Category = 4 //组
)

func init() {
	enums = make(map[Category]*enum)
	enums[Dept] = newEnum(int(Dept), "部门")
	enums[Company] = newEnum(int(Company), "公司")
	enums[Team] = newEnum(int(Team), "团队")
	enums[Group] = newEnum(int(Group), "组")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[Category]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := Category(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// Category 种类
type Category int

func (this Category) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this Category) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this Category) Code() int {
	return this.get().code
}

func (this Category) Message() string {
	return this.get().message
}
