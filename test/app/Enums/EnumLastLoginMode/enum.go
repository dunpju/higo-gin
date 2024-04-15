package EnumLastLoginMode

import (
	"fmt"
)

var (
	enums map[LastLoginMode]*enum
)

const (
	Mini LastLoginMode = 1 //小程序
	Pc LastLoginMode = 2 //pc端后台
	PcEnterpriseWechat LastLoginMode = 3 //pc端企业微信扫码
)

func init() {
	enums = make(map[LastLoginMode]*enum)
	enums[Mini] = newEnum(int(Mini), "小程序")
	enums[Pc] = newEnum(int(Pc), "pc端后台")
	enums[PcEnterpriseWechat] = newEnum(int(PcEnterpriseWechat), "pc端企业微信扫码")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[LastLoginMode]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := LastLoginMode(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// LastLoginMode 最后登录方式
type LastLoginMode int

func (this LastLoginMode) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this LastLoginMode) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this LastLoginMode) Code() int {
	return this.get().code
}

func (this LastLoginMode) Message() string {
	return this.get().message
}
