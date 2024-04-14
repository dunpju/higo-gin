package EnumWeixinIsrisky

import (
	"fmt"
)

var (
	enums map[WeixinIsrisky]*enum
)

const (
	Unknown WeixinIsrisky = 1 //未发起检测
	Nought WeixinIsrisky = 0 //暂未检测到风险
	Has WeixinIsrisky = 1 //风险
)

func init() {
	enums = make(map[WeixinIsrisky]*enum)
	enums[Unknown] = newEnum(int(Unknown), "未发起检测")
	enums[Nought] = newEnum(int(Nought), "暂未检测到风险")
	enums[Has] = newEnum(int(Has), "风险")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[WeixinIsrisky]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := WeixinIsrisky(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// WeixinIsrisky 微信检测结果
type WeixinIsrisky int

func (this WeixinIsrisky) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this WeixinIsrisky) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this WeixinIsrisky) Code() int {
	return this.get().code
}

func (this WeixinIsrisky) Message() string {
	return this.get().message
}
