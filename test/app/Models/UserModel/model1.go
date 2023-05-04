package UserModel

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-gin/higo/sql"
	"log"
)

func (this *Impl) UserById(id int, columns ...string) {
	fmt.Println(55, this.TableName())
	this.Mapper(squirrel.Select(columns...).From(this.TableName()).Where("id=?", id).ToSql()).Query().Scan(&this)
}

func (this *Impl) AddUser(uname string, tel string, score int) *higo.Orm {
	fmt.Println("tttt")
	fmt.Println(this.Builder().Select("uname", "u_tel").From("user_table").Where("id = ?", 9).ToSql())
	fmt.Println(this.Builder().Select("uname", "u_tel").From("user_table").Where("id = ?", 10).ToSql())
	fmt.Println(this.Builder().Select("uname", "u_tel").From("user_table").Where("id = ?", 11).ToSql())
	fmt.Println(this.Insert(this.TableName()).Set("uname", uname).Set("u_tel", tel).Set("score", score).ToSql())
	fmt.Println(sql.Insert(this.TableName()).InsertBuilder().Columns("uname", "u_tel", "score").Values(uname, tel, score).ToSql())
	fmt.Println("update")
	b := this.Update(this.TableName())
	b.Set("uname", "张三")
	b.Set("score", 5)
	b.Where("id", 4)
	fmt.Println(b.ToSql())
	fmt.Println(this.Update(this.TableName()).Set("uname", "张三").Set("score", 5).Where("id", 3).ToSql())
	fmt.Println(this.Update(this.TableName()).Set("uname", "张三").Set("score", 5).Where("id", 2).ToSql())
	fmt.Println(sql.Update(this.TableName()).Set("uname", "张三").Set("score", 5).Where("id", 1).ToSql())
	log.Println("ggggg")
	return this.Mapper(sql.Insert(this.TableName()).InsertBuilder().Columns("uname", "u_tel", "score").Values(uname, tel, score).ToSql())
}

func (this *Impl) Paginate1(perPage, page uint64) *higo.Pager {
	models := make([]*Impl, 0)
	pager := higo.NewPager(perPage, page)
	this.Table(this.TableName()).Where("uname like ?", "%werwerwer%").Paginate(pager).Find(&models)
	pager.Items = models
	fmt.Println(pager)
	fmt.Println(pager.Items)
	for _, v := range pager.Items.([]*Impl) {
		fmt.Println(v)
	}
	return pager
}

func (this *Impl) Add(uname string, tel string, score int) {
	//this.Paginate(2, 1)
	//u1 := this.New()
	//u := this.AddUser(uname, tel, score)
	//coinModel := CoinModel.New()
	//coin := CoinModel.New().AddCoin(uname, score)
	//user := this.New()
	//this.Find(user)
	//fmt.Println(user)
	//fmt.Println(this.TableName())
	//higo.Begin(u, coin).Transaction(func() error {
	//	higo.Result(u.Exec().Error).Unwrap()
	//	co := coin.Exec()
	//	higo.Result(co.Error).Unwrap()
	//	fmt.Println(1)
	//	fmt.Println(co.Value)
	//	coin.Last(&coinModel)
	//	higo.Result(u.Last(&coinModel).Error).Unwrap()
	//	fmt.Printf("%T\n", u1)
	//	higo.Result(u.Last(u1).Error).Unwrap()
	//	fmt.Println(coinModel)
	//	fmt.Println(u1)
	//	panic(fmt.Errorf("test"))
	//	return nil
	//})
}
