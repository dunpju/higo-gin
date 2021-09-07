package UserModel

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/higo/sql"
	"github.com/dengpju/higo-gin/test/app/Consts"
	"github.com/dengpju/higo-gin/test/app/Models/CoinModel"
	"github.com/dengpju/higo-ioc/injector"
	"log"
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

func (this *UserModelImpl) RegisterValidator() higo.Valid {
	//The custom tag
	//require import "gitee.com/dengpju/higo-code/code"
	// example
	return higo.RegisterValid(this).
		Tag("UserName",
			higo.Rule("required", Consts.Success),
			higo.Rule("min=5", Consts.Success)).
		Tag("Utel",
			higo.Rule("required", Consts.Success),
			higo.Rule("min=4", Consts.Success))
}

func (this *UserModelImpl) Exist() bool {
	return this.Id > 0
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
	fmt.Println("tttt")
	fmt.Println(this.Builder().Select("uname", "u_tel").From("user_table").Where("id = ?", 9).ToSql())
	fmt.Println(this.Builder().Select("uname", "u_tel").From("user_table").Where("id = ?", 10).ToSql())
	fmt.Println(this.Builder().Select("uname", "u_tel").From("user_table").Where("id = ?", 11).ToSql())
	fmt.Println(this.Insert(this.TableName()).
		Set("uname", uname).
		Set("u_tel", tel).
		Set("score", score).ToSql())
	fmt.Println(sql.Insert(this.TableName()).InsertBuilder().
		Columns("uname", "u_tel", "score").
		Values(uname, tel, score).
		ToSql())
	fmt.Println("update")
	b := this.Update(this.TableName())
	b.Set("uname", "张三")
	b.Set("score", 5)
	b.Where("id", 4)
	fmt.Println(b.ToSql())
	fmt.Println(this.Update(this.TableName()).
		Set("uname", "张三").
		Set("score", 5).
		Where("id", 3).ToSql())
	fmt.Println(this.Update(this.TableName()).
		Set("uname", "张三").
		Set("score", 5).
		Where("id", 2).
		ToSql())
	fmt.Println(sql.Update(this.TableName()).
		Set("uname", "张三").
		Set("score", 5).
		Where("id", 1).
		ToSql())

	log.Println("ggggg")
	return this.Mapper(sql.Insert(this.TableName()).InsertBuilder().
		Columns("uname", "u_tel", "score").
		Values(uname, tel, score).
		ToSql())
}

func (this *UserModelImpl) Paginate(perPage, page uint64) *higo.Pager {
	models := make([]*UserModelImpl, 0)
	pager := higo.NewPager(perPage, page)
	this.Table(this.TableName()).Where("uname like ?", "%werwerwer%").Paginate(pager, &models)
	pager.Items = models
	fmt.Println(pager)
	fmt.Println(pager.Items)
	for _, v := range pager.Items.([]*UserModelImpl) {
		fmt.Println(v)
	}
	return pager
}

func (this *UserModelImpl) Add(uname string, tel string, score int) {
	//higo.Result(this.AddUser(uname, tel, score).Exec().Error).Unwrap()
	this.Paginate(2, 1)
	//log.Fatalln("wanc")
	u1 := this.New()
	u := this.AddUser(uname, tel, score)

	coinModel := CoinModel.New()
	coin := CoinModel.New().AddCoin(uname, score)
	user := this.New()
	this.Find(user)
	fmt.Println(user)
	fmt.Println(this.TableName())
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
		panic(fmt.Errorf("test"))
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
