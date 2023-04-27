package AdminEntity

import (
	"github.com/dengpju/higo-gin/higo"
	"time"
)

type Impl struct {
	isEdit      bool
	currentFlag higo.Flag
	AdminId     int64       `json:"admin_id"`
	AdminName   string      `json:"admin_name"`
	UserId      int64       `json:"user_id"`
	State       int         `json:"state"`
	IsSuper     int         `json:"is_super"`
	Password    string      `json:"password"`
	CreateTime  time.Time   `json:"create_time" comment:"创建时间"`
	UpdateTime  time.Time   `json:"update_time" comment:"更新时间"`
	DeleteTime  interface{} `json:"delete_time" comment:"删除时间"`
}

func New() *Impl {
	t := time.Now()
	return &Impl{CreateTime: t, UpdateTime: t}
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
	return this.AdminId == 0
}
