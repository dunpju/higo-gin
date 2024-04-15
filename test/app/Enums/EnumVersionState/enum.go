package EnumVersionState

import (
	"fmt"
)

var (
	enums map[VersionState]*enum
)

const (
	Issue VersionState = 1 //发布
	Waiting VersionState = 2 //待发布
)

func init() {
	enums = make(map[VersionState]*enum)
	enums[Issue] = newEnum(int(Issue), "发布")
	enums[Waiting] = newEnum(int(Waiting), "待发布")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[VersionState]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := VersionState(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// VersionState 版本状态
type VersionState int

func (this VersionState) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this VersionState) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this VersionState) Code() int {
	return this.get().code
}

func (this VersionState) Message() string {
	return this.get().message
}
