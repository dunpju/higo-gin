package UserModel

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-gin/test/app/errcode"
	"github.com/dunpju/higo-ioc/injector"
)

type Impl struct {
	Id    int    `gorm:"column:id" json:"id" comment:""`
	Uname string `gorm:"column:uname" json:"uname" comment:""`
	UTel  string `gorm:"column:u_tel" json:"u_tel" comment:""`
	Score int    `gorm:"column:score" json:"score" comment:""`
}

// init Validator
func init() {
	New().RegisterValidator()
}

func New(attrs ...higo.Property) *Impl {
	impl := &Impl{}
	higo.Propertys(attrs).Apply(impl)
	injector.BeanFactory.Apply(impl)
	return impl
}

func (this *Impl) New() higo.IClass {
	return New()
}

func (this *Impl) TableName() string {
	return "ts_user"
}

func (this *Impl) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}

// The custom tag, binding the tag eg: binding:"custom_tag_name"
// require import "gitee.com/dengpju/higo-code/code"
//
// example code:
//
//	func (this *StructName) RegisterValidator() higo.Valid {
//		return higo.RegisterValid(this).
//			Tag("custom_tag_name",
//				higo.Rule("required", Codes.Success),
//				higo.Rule("min=5", Codes.Success))
//	 Or
//	 return higo.Verifier() // Manual call Register Validate: higo.Validate(verifier)
//	}
func (this *Impl) RegisterValidator() *higo.Verify {
	return higo.RegisterValidator(this).Tag("custom_tag_name", higo.Rule("required", errcode.EnumError), higo.Rule("min=5", errcode.EnumError))
}

func (this *Impl) Exist() bool {
	return this.Id > 0
}
