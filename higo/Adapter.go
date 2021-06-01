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
	this.sql = sql
	this.args = args
	return this
}

func (this *Orm) Query() *gorm.DB {
	return this.Raw(this.sql, this.args)
}

func (this *Orm) Execute() *gorm.DB {
	return this.Exec(this.sql, this.args)
}
