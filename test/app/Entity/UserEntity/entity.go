package UserEntity

import (
	"github.com/dunpju/higo-orm/arm"
)

const (
	FlagDelete arm.Flag = iota + 1
	FlagUpdate
)

type Entity struct {
	_edit bool
	_flag arm.Flag
	Id    int    `gorm:"column:id" json:"id" comment:""`
	Uname string `gorm:"column:uname" json:"uname" comment:""`
	UTel  string `gorm:"column:u_tel" json:"u_tel" comment:""`
	Score int    `gorm:"column:score" json:"score" comment:""`
}

func New() *Entity {
	return &Entity{}

}

func (this *Entity) IsEdit() bool {
	return this._edit
}

func (this *Entity) Edit(isEdit bool) {
	this._edit = isEdit
}

func (this *Entity) Flag(flag arm.Flag) {
	this._flag = flag
	this._edit = true
}

func (this *Entity) Equals(flag arm.Flag) bool {
	return this._flag == flag
}

func (this *Entity) PrimaryEmpty() bool {
	return this.Id == 0
}
