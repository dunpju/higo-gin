package NewsModel

import (
	"github.com/Masterminds/squirrel"
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-ioc/injector"
	"github.com/jinzhu/gorm"
	"strings"
)

type Impl struct {
	*higo.Orm  `inject:"Bean.NewOrm()"`
	NewsId     int         `gorm:"column:news_id" json:"news_id" comment:"主键"`
	Title      string      `gorm:"column:title" json:"title" comment:"标题"`
	Clicknum   int         `gorm:"column:clicknum" json:"clicknum" comment:"点击量"`
	CreateTime interface{} `gorm:"column:create_time" json:"create_time" comment:"创建时间"`
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
	return "ts_news"
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
func (this *Impl) RegisterValidator() *higo.Valid {
	return higo.RegisterValid(this)
}

func (this *Impl) Exist() bool {
	return this.NewsId > 0
}

func (this *Impl) GetByNewsId(newsId int, columns ...string) *gorm.DB {
	return this.Mapper(squirrel.Select(columns...).From(this.TableName()).Where(NewsId+" = ?", newsId).ToSql()).Query()
}

func (this *Impl) GetByNewsIds(newsIds []string, columns ...string) *gorm.DB {
	return this.Mapper(squirrel.Select(columns...).From(this.TableName()).Where(NewsId+" IN(?)", strings.Join(newsIds, ",")).ToSql()).Query()
}

func (this *Impl) Paginate(perPage, page uint64) *higo.Pager {
	models := make([]*Impl, 0)
	pager := higo.NewPager(perPage, page)
	this.Table(this.TableName()).Paginate(pager).Find(&models)
	pager.Items = models
	return pager
}
