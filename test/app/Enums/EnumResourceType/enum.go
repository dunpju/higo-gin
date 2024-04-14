package EnumResourceType

import (
	"fmt"
)

var (
	enums map[ResourceType]*enum
)

const (
	Audio ResourceType = 1 //音频
	Video ResourceType = 2 //视频
	Image ResourceType = 3 //图片
)

func init() {
	enums = make(map[ResourceType]*enum)
	enums[Audio] = newEnum(int(Audio), "音频")
	enums[Video] = newEnum(int(Video), "视频")
	enums[Image] = newEnum(int(Image), "图片")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[ResourceType]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := ResourceType(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// ResourceType 资源类型
type ResourceType int

func (this ResourceType) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this ResourceType) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this ResourceType) Code() int {
	return this.get().code
}

func (this ResourceType) Message() string {
	return this.get().message
}
