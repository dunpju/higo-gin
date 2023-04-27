package DemoModel

type DemoModelImpl struct {
	Id int `gorm:"column:id"`
	Uname string `gorm:"column:uname"`
	Utel string `gorm:"column:u_tel"`
}

func New() *DemoModelImpl {
	return &DemoModelImpl{}
}

