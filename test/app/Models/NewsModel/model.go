package NewsModel

import (
	"github.com/dunpju/higo-orm/arm"
	"github.com/dunpju/higo-orm/him"
)

type Model struct {
	*arm.Model
	NewsId     int         `gorm:"column:news_id" json:"news_id" comment:"主键"`
	Title      string      `gorm:"column:title" json:"title" comment:"标题"`
	Clicknum   int         `gorm:"column:clicknum" json:"clicknum" comment:"点击量"`
	CreateTime interface{} `gorm:"column:create_time" json:"create_time" comment:"创建时间"`
}

func New(properties ...him.IProperty) *Model {
	return (&Model{}).New(properties...)
}

func TableName() *arm.TableName {
	return arm.NewTableName("ts_news")
}

func (this *Model) New(properties ...him.IProperty) *Model {
	err := arm.Connect(this)
	if err != nil {
		panic(err)
	}
	this.Property(properties...)
	return this
}

func (this *Model) Mutate(properties ...him.IProperty) arm.IModel {
	return New(properties...)
}

func (this *Model) Connection() string {
	return him.DefaultConnect
}

func (this *Model) TableName() *arm.TableName {
	return TableName()
}

func (this *Model) Apply(model *arm.Model) {
	this.Model = model
}

func (this *Model) Exist() bool {
	return this.NewsId > 0
}
