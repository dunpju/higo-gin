package NewsEntity

import (
	"github.com/dunpju/higo-gin/higo"
	"time"
)

type Impl struct {
	isEdit      bool
	currentFlag higo.Flag
	NewsId      int         `gorm:"column:news_id" json:"news_id" comment:"主键"`
	Title       string      `gorm:"column:title" json:"title" comment:"标题"`
	Clicknum    int         `gorm:"column:clicknum" json:"clicknum" comment:"点击量"`
	CreateTime  interface{} `gorm:"column:create_time" json:"create_time" comment:"创建时间"`
}

func New() *Impl {
	t := time.Now()
	return &Impl{CreateTime: tn}

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
	return this.NewsId == 0
}
