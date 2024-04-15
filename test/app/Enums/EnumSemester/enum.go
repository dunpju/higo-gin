package EnumSemester

import (
	"fmt"
)

var (
	enums map[Semester]*enum
)

const (
	Up Semester = 1 //上学期
	Down Semester = 2 //下学期
)

func init() {
	enums = make(map[Semester]*enum)
	enums[Up] = newEnum(int(Up), "上学期")
	enums[Down] = newEnum(int(Down), "下学期")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[Semester]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := Semester(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// Semester 学期
type Semester int

func (this Semester) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this Semester) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this Semester) Code() int {
	return this.get().code
}

func (this Semester) Message() string {
	return this.get().message
}
