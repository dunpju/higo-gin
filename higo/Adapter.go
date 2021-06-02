package higo

import (
	"fmt"
	"github.com/dengpju/higo-config/config"
	"github.com/dengpju/higo-logger/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"sync"
)

var (
	orm      *Orm
	onceGorm sync.Once
)

type Orm struct {
	*gorm.DB
	sql  string
	args []interface{}
}

func NewOrm() *Orm {
	onceGorm.Do(func() {
		confDefault := config.Db("DB.DEFAULT").(*config.Configure)
		args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			confDefault.Get("USERNAME").(string),
			confDefault.Get("PASSWORD").(string),
			confDefault.Get("HOST").(string),
			confDefault.Get("PORT").(string),
			confDefault.Get("DATABASE").(string),
			confDefault.Get("CHARSET").(string),
		)
		db, err := gorm.Open(confDefault.Get("DRIVER").(string), args)
		if err != nil {
			log.Fatal(err)
		}
		logger.Logrus.Infoln(fmt.Sprintf("DB %s:%s Connection success!", confDefault.Get("HOST").(string),
			confDefault.Get("PORT").(string)))
		db.SingularTable(true)
		db.DB().SetMaxIdleConns(5)
		db.DB().SetMaxOpenConns(10)
		orm = &Orm{DB: db}
	})
	return orm
}

func (this *Orm) Mapper(sql string, args []interface{}, err error) *Orm {
	if err != nil {
		panic(err.Error())
	}
	cloneDB := &Orm{}
	cloneDB.DB = orm.DB
	cloneDB.sql = sql
	cloneDB.args = args
	return cloneDB
}

func (this *Orm) setDB(db *gorm.DB) {
	this.DB = db
}

func (this *Orm) Query() *gorm.DB {
	if this.DB != nil {
		return this.DB.Exec(this.sql, this.args...)
	}
	return this.Raw(this.sql, this.args...)
}

func (this *Orm) Execute() *gorm.DB {
	if this.DB != nil {
		return this.DB.Exec(this.sql, this.args...)
	}
	return this.Exec(this.sql, this.args...)
}

func (this *Orm) BeginTransaction(orms ...*Orm) error {
	tx := this.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	if tx.Error != nil {
		panic(tx.Error)
	}
	for _, orm := range orms {
		if orm.Execute().Error != nil {
			tx.Rollback()
			panic(orm.Error)
		}
	}
	tx.Commit()
	return nil
}

type Orms []*Orm

func (this Orms) apply(tx *gorm.DB) {
	for _, sql := range this {
		sql.setDB(tx)
	}
}

//执行事务
func (this Orms) Execute(fn func() error) error {
	return orm.Transaction(func(tx *gorm.DB) error {
		this.apply(tx)
		return fn()
	})
}

func Mappers(orms ...*Orm) (list Orms) {
	list = orms
	return
}
