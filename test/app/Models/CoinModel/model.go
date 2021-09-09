package CoinModel

import (
	"github.com/Masterminds/squirrel"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-ioc/injector"
)

type Impl struct {
	*higo.Orm `inject:"Bean.NewOrm()"`
	Id        int    `gorm:"column:id"`
	Uname     string `gorm:"column:uname"`
	Coin      int    `gorm:"column:coin"`
}

func New(attrs ...higo.Property) *Impl {
	model := &Impl{}
	higo.Propertys(attrs).Apply(model)
	injector.BeanFactory.Apply(model)
	return model
}

func (this *Impl) New() higo.IClass {
	return New()
}

func (this *Impl) TableName() string {
	return "ts_coin"
}

func (this *Impl) AddCoin(uname string, coin int) *higo.Orm {
	return higo.Mapper(squirrel.Insert(this.TableName()).Columns("uname", "coin").Values(uname, coin).ToSql()).Orm
}
