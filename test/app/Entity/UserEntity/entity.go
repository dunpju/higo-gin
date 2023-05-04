package UserEntity

import (
	"github.com/dunpju/higo-gin/higo"
)

type Impl struct {
	isEdit      bool
	currentFlag higo.Flag
	Id          int    `gorm:"column:id" json:"id" comment:""`
	Uname       string `gorm:"column:uname" json:"uname" comment:""`
	UTel        string `gorm:"column:u_tel" json:"u_tel" comment:""`
	Score       int    `gorm:"column:score" json:"score" comment:""`
}

func New() *Impl {
	return &Impl{}

}

func (this *Impl) IsEdit() bool {
	return this.isEdit
}

func (this *Impl) SetIsEdit(isEdit bool) {
	this.isEdit = isEdit
}

func (this *Impl) SetFlag(flag higo.Flag) {
	this.currentFlag = flag
	this.isEdit = true
}

func (this *Impl) Flag() higo.Flag {
	return this.currentFlag
}

func (this *Impl) PriEmpty() bool {
	return this.Id == 0
}
