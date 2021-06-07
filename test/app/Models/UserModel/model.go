package UserModel

import (
	"fmt"
	"gitee.com/dengpju/higo-code/code"
	"github.com/Masterminds/squirrel"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Models/CoinModel"
	"github.com/dengpju/higo-ioc/injector"
)

type UserModelImpl struct {
	*higo.Orm `inject:"Bean.NewOrm()"`
	Id        int    `gorm:"column:id" json:"id" binding:"required"`
	Uname     string `gorm:"column:uname" json:"uname" binding:"UserName"`
	Utel      string `gorm:"column:u_tel" json:"utel" binding:"Utel"`
}

func init() {
	//初始化校验器
	New().RegisterValidator()
}

func New(attrs ...higo.Property) *UserModelImpl {
	u := &UserModelImpl{}
	higo.Propertys(attrs).Apply(u)
	injector.BeanFactory.Apply(u)
	return u
}

func (this *UserModelImpl) New() higo.IClass {
	return New()
}

func (this *UserModelImpl) TableName() string {
	return "ts_user"
}

func (this *UserModelImpl) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}

func (this *UserModelImpl) RegisterValidator() {
	higo.RegisterValid(this).
		Tag("UserName",
			higo.Rule("required", code.Message("20000@UserName必须填1")),
			higo.Rule("min=5", code.Message("20000@UserName大于5"))).
		Tag("Utel",
			higo.Rule("required", code.Message("20000@Utel必须填")),
			higo.Rule("min=4", code.Message("20000@Utel大于4")))
}

func (this *UserModelImpl) UserById(id int, columns ...string) {
	fmt.Println(55, this.TableName())
	this.Mapper(squirrel.
		Select(columns...).
		From(this.TableName()).
		Where("id=?", id).
		ToSql()).
		Query().
		Scan(&this)
}

func (this *UserModelImpl) AddUser(uname string, tel string, score int) *higo.Orm {
	return this.Mapper(squirrel.Insert(this.TableName()).
		Columns("uname", "u_tel", "score").
		Values(uname, tel, score).
		ToSql())
}

func (this *UserModelImpl) Add(uname string, tel string, score int) {
	u1 := this.New()
	u := this.AddUser(uname, tel, score)
	coinModel := CoinModel.New()
	coin := CoinModel.New().AddCoin(uname, score)
	//方法一:

	higo.Begin(u, coin).Transaction(func() error {
		higo.Result(u.Exec().Error).Unwrap()
		co := coin.Exec()
		higo.Result(co.Error).Unwrap()
		fmt.Println(1)
		fmt.Println(co.Value)
		coin.Last(&coinModel)
		higo.Result(u.Last(&coinModel).Error).Unwrap()
		fmt.Printf("%T\n", u1)
		higo.Result(u.Last(u1).Error).Unwrap()
		fmt.Println(coinModel)
		fmt.Println(u1)
		return nil
	})

	//方法二:
	/**
	this.Begin(u, coin).Transaction(func() error {
		higo.Result(u.Exec().Error).Unwrap()
		higo.Result(coin.Exec().Error).Unwrap()
		fmt.Println(1)
		coin.Last(&coinModel)
		higo.Result(u.Last(&coinModel).Error).Unwrap()
		fmt.Printf("%T\n", u1)
		higo.Result(u.Last(u1).Error).Unwrap()
		fmt.Println(coinModel)
		fmt.Println(u1)
		return nil
	})

	*/
}
