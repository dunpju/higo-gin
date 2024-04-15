package EnumVersionIsUsed

import (
	"fmt"
)

var (
	enums map[VersionIsUsed]*enum
)

const (
	Yes VersionIsUsed = 1 //是
	No VersionIsUsed = 2 //否
)

func init() {
	enums = make(map[VersionIsUsed]*enum)
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

func Enums() map[VersionIsUsed]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := VersionIsUsed(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// VersionIsUsed 是否正在使用
type VersionIsUsed int

func (this VersionIsUsed) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this VersionIsUsed) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this VersionIsUsed) Code() int {
	return this.get().code
}

func (this VersionIsUsed) Message() string {
	return this.get().message
}
