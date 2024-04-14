package EnumResourceMaterialType

import (
	"fmt"
)

var (
	enums map[ResourceMaterialType]*enum
)

const (
	Word ResourceMaterialType = 1 //word
	Pdf ResourceMaterialType = 2 //pdf
)

func init() {
	enums = make(map[ResourceMaterialType]*enum)
	enums[Word] = newEnum(int(Word), "word")
	enums[Pdf] = newEnum(int(Pdf), "pdf")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[ResourceMaterialType]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := ResourceMaterialType(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// ResourceMaterialType 资源材料类型
type ResourceMaterialType int

func (this ResourceMaterialType) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this ResourceMaterialType) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this ResourceMaterialType) Code() int {
	return this.get().code
}

func (this ResourceMaterialType) Message() string {
	return this.get().message
}
