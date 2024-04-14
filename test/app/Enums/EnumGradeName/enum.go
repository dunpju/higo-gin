package EnumGradeName

import (
	"fmt"
)

var (
	enums map[GradeName]*enum
)

const (
	KD GradeName = "KD" //婴班
	KC GradeName = "KC" //小班
	KB GradeName = "KB" //中班
	KA GradeName = "KA" //大班
)

func init() {
	enums = make(map[GradeName]*enum)
	enums[KD] = newEnum(string(KD), "婴班")
	enums[KC] = newEnum(string(KC), "小班")
	enums[KB] = newEnum(string(KB), "中班")
	enums[KA] = newEnum(string(KA), "大班")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[GradeName]*enum {
	return enums
}

func Inspect(value string) error {
	_, err := GradeName(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// GradeName 幼儿园年级
type GradeName string

func (this GradeName) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this GradeName) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this GradeName) Code() int {
	return this.get().code
}

func (this GradeName) Message() string {
	return this.get().message
}
