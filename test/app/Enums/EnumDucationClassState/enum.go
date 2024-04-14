package EnumDucationClassState

import (
	"fmt"
)

var (
	enums map[DucationClassState]*enum
)

const (
	Unknown DucationClassState = 0 //未开始
	Ing DucationClassState = 1 //进行中
	Archive DucationClassState = 2 //已归档
)

func init() {
	enums = make(map[DucationClassState]*enum)
	enums[Unknown] = newEnum(int(Unknown), "未开始")
	enums[Ing] = newEnum(int(Ing), "进行中")
	enums[Archive] = newEnum(int(Archive), "已归档")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[DucationClassState]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := DucationClassState(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// DucationClassState 状态
type DucationClassState int

func (this DucationClassState) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this DucationClassState) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this DucationClassState) Code() int {
	return this.get().code
}

func (this DucationClassState) Message() string {
	return this.get().message
}
