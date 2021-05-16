package UserModel

import "github.com/dengpju/higo-gin/higo"

type UserModelImpl struct {
	Id    int    `gorm:"column:id"`
	Uname string `gorm:"column:uname"`
	Utel  string `gorm:"column:u_tel"`
}

func New(attrs ...higo.Property) *UserModelImpl {
	u := &UserModelImpl{}
	higo.Propertys(attrs).Apply(u)
	return u
}

func (this *UserModelImpl) New() higo.IClass {
	return New()
}
