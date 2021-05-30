package UserModel

import (
	"github.com/dengpju/higo-gin/higo"
)

type UserModelImpl struct {
	Id    int    `gorm:"column:id" json:"id" binding:"required"`
	Utel  string `gorm:"column:u_tel" json:"utel" binding:"Utel"`
	Uname string `gorm:"column:uname" json:"uname" binding:"UserName"`
}

func New(attrs ...higo.Property) *UserModelImpl {
	u := &UserModelImpl{}
	higo.Propertys(attrs).Apply(u)
	return u
}

func (this *UserModelImpl) New() higo.IClass {
	return New()
}

func (this *UserModelImpl) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}

func (this *UserModelImpl) InitValidator() higo.Valid {
	return higo.RegisterValid(this).
		Tag("UserName",
			higo.Rule("required", "20000@UserName必须填"),
			higo.Rule("min=5", "20000@UserName必须填大于5")).
		Tag("Utel",
			higo.Rule("required", "20000@Utel必须填"),
			higo.Rule("min=4", "20000@Utel大于4"))
}
