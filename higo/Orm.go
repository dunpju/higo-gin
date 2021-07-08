package higo

import (
	"fmt"
	"github.com/dengpju/higo-config/config"
	"github.com/dengpju/higo-logger/logger"
	"github.com/dengpju/higo-throw/exception"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"sync"
)

var (
	orm          *Orm
	onceGorm     sync.Once
	dbConfigOnce sync.Once
	confDefault  *config.Configure
	dbConfig     *Dbconfig
)

type Dbconfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Charset  string
	Driver   string
	Prefix   string
}

type Orm struct {
	*gorm.DB
	sql  string
	args []interface{}
	orms []*Orm
}

func GetDbConfig() *Dbconfig {
	return dbConfig
}

func (this *Orm) Args() []interface{} {
	return this.args
}

func (this *Orm) Sql() string {
	return this.sql
}

func newGorm() *gorm.DB {
	dbConfigOnce.Do(func() {
		confDefault = config.Db("DB.DEFAULT").(*config.Configure)
		dbConfig = &Dbconfig{Username: confDefault.Get("USERNAME").(string),
			Password: confDefault.Get("PASSWORD").(string),
			Host:     confDefault.Get("HOST").(string),
			Port:     confDefault.Get("PORT").(string),
			Database: confDefault.Get("DATABASE").(string),
			Charset:  confDefault.Get("CHARSET").(string),
			Driver:   confDefault.Get("DRIVER").(string),
			Prefix:   confDefault.Get("PREFIX").(string),
		}
	})
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		dbConfig.Charset,
	)
	db, err := gorm.Open(dbConfig.Driver, args)
	if err != nil {
		log.Fatal(err)
	}
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
		logger.Logrus.Infoln(fmt.Sprintf("DB %s:%s Connection success!", dbConfig.Host,
			dbConfig.Port))
	})
	return orm
}

func MultiOrm() *Orm {
	return newOrm()
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
	return this.DB.Raw(this.sql, this.args...)
}

func (this *Orm) Exec() *gorm.DB {
	return this.DB.Exec(this.sql, this.args...)
}

func (this *Orm) Begin(orms ...*Orm) *Orm {
	this.orms = orms
	return this
}

func (this *Orm) apply(tx *gorm.DB) {
	for _, sql := range this.orms {
		sql.setDB(tx)
	}
}

func (this *Orm) Transaction(fn func() error) {
	err := this.DB.Transaction(func(tx *gorm.DB) (err error) {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				err = r.(error)
				return
			}
		}()
		this.apply(tx)
		err = fn()
		return
	})
	if err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			exception.Throw(exception.Message(e.Message),
				exception.Code(int(e.Number)),
				exception.Data(nil))
		} else {
			panic(err)
		}
	}
}
