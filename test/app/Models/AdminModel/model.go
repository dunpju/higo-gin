package AdminModel

import (
	"github.com/dunpju/higo-orm/arm"
	"github.com/dunpju/higo-orm/him"
	"time"
)

type Model struct {
	*arm.Model
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

func New(properties ...him.IProperty) *Model {
	return (&Model{}).New(properties...)
}

func TableName() *arm.TableName {
	return arm.NewTableName("tl_admin")
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
	return this.AdminId > 0
}
