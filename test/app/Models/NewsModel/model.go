package NewsModel

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-ioc/injector"
)

type Impl struct {
	NewsId     int         `gorm:"column:news_id" json:"news_id" comment:"主键"`
	Title      string      `gorm:"column:title" json:"title" comment:"标题"`
	Clicknum   int         `gorm:"column:clicknum" json:"clicknum" comment:"点击量"`
	CreateTime interface{} `gorm:"column:create_time" json:"create_time" comment:"创建时间"`
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

func (this *Impl) Exist() bool {
	return this.NewsId > 0
}
