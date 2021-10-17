package UserModel

import (
	//"gitee.com/dengpju/higo-code/code"
	"github.com/Masterminds/squirrel"
    "github.com/dengpju/higo-gin/higo"
    "github.com/dengpju/higo-ioc/injector"
    "github.com/jinzhu/gorm"
    "strings"
)

type Impl struct {
	*higo.Orm    `inject:"Bean.NewOrm()"`
	Id    int    `gorm:"column:id" json:"id" comment:""`
	Uname    string    `gorm:"column:uname" json:"uname" comment:""`
	UTel    string    `gorm:"column:u_tel" json:"u_tel" comment:""`
	Score    int    `gorm:"column:score" json:"score" comment:""`
}

//init Validator
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

//The custom tag, binding the tag eg: binding:"custom_tag_name"
//require import "gitee.com/dengpju/higo-code/code"
//
//example code:
//func (this *StructName) RegisterValidator() higo.Valid {
//	return higo.RegisterValid(this).
//		Tag("custom_tag_name",
//			higo.Rule("required", Codes.Success),
//			higo.Rule("min=5", Codes.Success))
//  Or
//  return higo.Verifier() // Manual call Register Validate: higo.Validate(verifier)
//}
func (this *Impl) RegisterValidator() *higo.Verify {
    return higo.RegisterValidator(this)
}

func (this *Impl) Exist() bool {
	return this.Id > 0
}

func (this *Impl) GetById(id int, columns ...string) *gorm.DB {
	return this.Mapper(squirrel.Select(columns...).From(this.TableName()).Where(Id + " = ?", id).ToSql()).Query()
}

func (this *Impl) GetByIds(ids []string, columns ...string) *gorm.DB {
	return this.Mapper(squirrel.Select(columns...).From(this.TableName()).Where(Id + " IN(?)", strings.Join(ids, ",")).ToSql()).Query()
}

func (this *Impl) Paginate(perPage, page uint64) *higo.Pager {
	models := make([]*Impl, 0)
	pager := higo.NewPager(perPage, page)
	this.Table(this.TableName()). /**Where().*/ Paginate(pager).Find(&models)
	pager.Items = models
	return pager
}