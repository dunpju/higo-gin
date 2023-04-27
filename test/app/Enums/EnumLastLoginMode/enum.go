package EnumLastLoginMode

import "github.com/dengpju/higo-enum/enum"

var e LastLoginMode

func Inspect(value int) error {
	return e.Inspect(value)
}

//最后登录方式
type LastLoginMode int

func (this LastLoginMode) Name() string {
	return "LastLoginMode"
}

func (this LastLoginMode) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this LastLoginMode) Message() string {
	return enum.String(this)
}

const (
	Mini LastLoginMode = 1 //小程序
	Pc LastLoginMode = 2 //pc端后台
	PcEnterpriseWechat LastLoginMode = 3 //pc端企业微信扫码
)

func (this LastLoginMode) Register() enum.Message {
	return make(enum.Message).
	    Put(Mini, "小程序").
	    Put(Pc, "pc端后台").
	    Put(PcEnterpriseWechat, "pc端企业微信扫码")
}