package EnumIsRecheck

import (
	"fmt"
)

var (
	enums map[IsRecheck]*enum
)

const (
	Yes IsRecheck = 1 //是
	No IsRecheck = 2 //否
)

func init() {
	enums = make(map[IsRecheck]*enum)
	enums[Yes] = newEnum(int(Yes), "是")
	enums[No] = newEnum(int(No), "否")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[IsRecheck]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := IsRecheck(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// IsRecheck 是否复核
type IsRecheck int

func (this IsRecheck) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this IsRecheck) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this IsRecheck) Code() int {
	return this.get().code
}

func (this IsRecheck) Message() string {
	return this.get().message
}
