package User

import (
	"github.com/dunpju/higo-orm/arm"
	"github.com/dunpju/higo-orm/him"
)

const (
	Id    arm.Fields = "id"
	Uname arm.Fields = "uname"
	UTel  arm.Fields = "u_tel"
	Score arm.Fields = "score"
)

type Model struct {
	*arm.Model
	Id    int    `gorm:"column:id" json:"id" comment:""`
	Uname string `gorm:"column:uname" json:"uname" comment:""`
	UTel  string `gorm:"column:u_tel" json:"u_tel" comment:""`
	Score int    `gorm:"column:score" json:"score" comment:""`
}

func New(properties ...him.IProperty) *Model {
	return (&Model{}).New(properties...)
}

func TableName() *arm.TableName {
	return arm.NewTableName("ts_user")
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
	return this.Id > 0
}

// WithId
func WithId(id int) him.IProperty {
	return him.SetProperty(func(obj any) {
		obj.(*Model).Id = id
	})
}
