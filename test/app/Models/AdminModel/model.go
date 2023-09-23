package AdminModel

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-ioc/injector"
	"time"
)

type Impl struct {
	AdminId    int64       `gorm:"column:admin_id" json:"admin_id" comment:"主键"`
	AdminName  string      `gorm:"column:admin_name" json:"admin_name" comment:"管理员名称"`
	UserId     int64       `gorm:"column:user_id" json:"user_id" comment:"用户id:对应user表"`
	State      int         `gorm:"column:state" json:"state" comment:"状态:1-启用,2-禁用"`
	IsSuper    int         `gorm:"column:is_super" json:"is_super" comment:"是否超级管理员:1-是,2-否"`
	Password   string      `gorm:"column:password" json:"password" comment:"密码"`
	CreateTime time.Time   `gorm:"column:create_time" json:"create_time" comment:"创建时间"`
	UpdateTime time.Time   `gorm:"column:update_time" json:"update_time" comment:"更新时间"`
	DeleteTime interface{} `gorm:"column:delete_time" json:"delete_time" comment:"删除时间"`
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
	return "tl_admin"
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
//	func (this *ModelImpl) RegisterValidator() higo.Valid {
//		return higo.RegisterValid(this).
//			Tag("custom_tag_name",
//				higo.Rule("required", Codes.Success),
//				higo.Rule("min=5", Codes.Success))
//	 Or
//	 return higo.Verifier() // Manual call Register Validate: higo.Validate(verifier)
//	}
func (this *Impl) RegisterValidator() *higo.Verify {
	return higo.RegisterValidator(this)
}

func (this *Impl) Exist() bool {
	return this.AdminId > 0
}
