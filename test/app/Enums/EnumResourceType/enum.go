package EnumResourceType

import "github.com/dunpju/higo-enum/enum"

var e ResourceType

func Inspect(value int) error {
	return e.Inspect(value)
}

//资源类型
type ResourceType int

func (this ResourceType) Name() string {
	return "ResourceType"
}

func (this ResourceType) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this ResourceType) Message() string {
	return enum.String(this)
}

const (
	Audio ResourceType = 1 //音频
	Video ResourceType = 2 //视频
	Image ResourceType = 3 //图片
)

func (this ResourceType) Register() enum.Message {
	return make(enum.Message).
		Put(Audio, "音频").
		Put(Video, "视频").
		Put(Image, "图片")
}
