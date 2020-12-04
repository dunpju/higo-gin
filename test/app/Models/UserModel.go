package Models

type UserModel struct {
	Id int `gorm:"column:id"`
	Uname string `gorm:"column:uname"`
	Utel string `gorm:"column:u_tel"`
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (this *UserModel)String() string {
	return "UserModel"
}