package EnumGender

import (
	"fmt"
)

var (
	enums map[Gender]*enum
)

const (
	Unknown Gender = 0 //未知
	Male Gender = 1 //男
	Female Gender = 2 //女
)

func init() {
	enums = make(map[Gender]*enum)
	enums[Unknown] = newEnum(int(Unknown), "未知")
	enums[Male] = newEnum(int(Male), "男")
	enums[Female] = newEnum(int(Female), "女")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[Gender]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := Gender(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// Gender 是否正在使用
type Gender int

func (this Gender) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this Gender) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this Gender) Code() int {
	return this.get().code
}

func (this Gender) Message() string {
	return this.get().message
}
