package EnumWeixinIsrisky

import "github.com/dengpju/higo-enum/enum"

var e WeixinIsrisky

func Inspect(value int) error {
	return e.Inspect(value)
}

//微信检测结果
type WeixinIsrisky int

func (this WeixinIsrisky) Name() string {
	return "WeixinIsrisky"
}

func (this WeixinIsrisky) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this WeixinIsrisky) Message() string {
	return enum.String(this)
}

const (
	Unknown WeixinIsrisky = 1 //未发起检测
	Nought WeixinIsrisky = 0 //暂未检测到风险
	Has WeixinIsrisky = 1 //风险
)

func (this WeixinIsrisky) Register() enum.Message {
	return make(enum.Message).
	    Put(Unknown, "未发起检测").
	    Put(Nought, "暂未检测到风险").
	    Put(Has, "风险")
}