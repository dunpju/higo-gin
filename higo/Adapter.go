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
	orms []*Orm
}

func newGorm() *gorm.DB {
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
	return db
}

func newOrm() *Orm {
	return &Orm{DB: newGorm(), orms: make([]*Orm, 0)}
}

func NewOrm() *Orm {
	onceGorm.Do(func() {
		orm = newOrm()
	})
	return orm
}

func (this *Orm) Mapper(sql string, args []interface{}, err error) *Orm {
	if err != nil {
		panic(err.Error())
	}
	cloneDB := newOrm()
	cloneDB.DB = orm.DB
	cloneDB.sql = sql
	cloneDB.args = args
	return cloneDB
}

func (this *Orm) setDB(db *gorm.DB) {
	this.DB = db
}

func (this *Orm) Query() *gorm.DB {
	return this.Raw(this.sql, this.args...)
}

func (this *Orm) Execute() *gorm.DB {
	if this.DB != nil {
		return this.DB.Exec(this.sql, this.args...)
	}
	return this.Exec(this.sql, this.args...)
}

func (this *Orm) Begin(orms ...*Orm) *gorm.DB {
	this.orms = orms
	this.apply(this.DB)
	return this.DB
}

func (this *Orm) apply(tx *gorm.DB) {
	for _, sql := range this.orms {
		sql.setDB(tx)
	}
}

func Mapper(sql string, args []interface{}, err error) *Orm {
	if err != nil {
		panic(err.Error())
	}
	cloneDB := newOrm()
	cloneDB.sql = sql
	cloneDB.args = args
	return cloneDB
}
