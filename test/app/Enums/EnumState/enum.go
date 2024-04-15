package EnumState

import (
	"fmt"
)

var (
	enums map[State]*enum
)

const (
	Issue State = 1 //发布
	Draft State = 2 //草稿
)

func init() {
	enums = make(map[State]*enum)
	enums[Issue] = newEnum(int(Issue), "发布")
	enums[Draft] = newEnum(int(Draft), "草稿")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[State]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := State(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// State 状态
type State int

func (this State) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this State) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this State) Code() int {
	return this.get().code
}

func (this State) Message() string {
	return this.get().message
}
