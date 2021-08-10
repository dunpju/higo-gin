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
	"math"
	"reflect"
	"strings"
	"sync"
)

var (
	orm          *Orm
	onceGorm     sync.Once
	dbConfigOnce sync.Once
	confDefault  *config.Configure
	dbConfig     *Dbconfig
	logMode      bool
	maxIdle      int
	maxOpen      int
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
	sql   string
	args  []interface{}
	orms  []*Orm
	table string
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

func newGorm(mapper bool) *gorm.DB {
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
		logMode = confDefault.Get("LOG_MODE").(bool)
		maxIdle = confDefault.Get("MAX_IDLE").(int)
		maxOpen = confDefault.Get("MAX_OPEN").(int)
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
	db.LogMode(logMode)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxOpen)
	if mapper == false {
		if db.Callback().Query().Get("gorm:Query") == nil {
			db.Callback().Query().Before("gorm:Query").Register("Query", sqlReplace)
		}
		if db.Callback().RowQuery().Get("gorm:Query") == nil {
			db.Callback().RowQuery().Before("gorm:RowQuery").Register("RowQuery", sqlReplace)
		}
		if db.Callback().Create().Get("gorm:Create") == nil {
			db.Callback().Create().Before("gorm:Create").Register("Create", sqlReplace)
		}
		if db.Callback().Update().Get("gorm:Update") == nil {
			db.Callback().Update().Before("gorm:Update").Register("Update", sqlReplace)
		}
		if db.Callback().Delete().Get("gorm:Delete") == nil {
			db.Callback().Delete().Before("gorm:Delete").Register("Delete", sqlReplace)
		}
	}
	return db
}

func sqlReplace(scope *gorm.Scope) {
	sql := scope.SQL
	s := reflect.ValueOf(scope.SQLVars)
	if orm, ok := scope.Value.(*Orm); ok {
		sql = orm.sql
		s = reflect.ValueOf(orm.args)
	}
	for i := 0; i < s.Len(); i++ {
		sql = strings.Replace(sql, "?", "'%v'", 1)
		sql = fmt.Sprintf(sql, s.Index(i))
	}
	if logMode {
		logger.Logrus.Debugln(sql)
	}
}

func newOrm() *Orm {
	return &Orm{DB: newGorm(false), orms: make([]*Orm, 0)}
}

func mapperOrm() *Orm {
	return &Orm{DB: newGorm(true), orms: make([]*Orm, 0)}
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
	return mapperOrm()
}

func (this *Orm) Mapper(sql string, args []interface{}, err error) *Orm {
	if err != nil {
		panic(err.Error())
	}
	clone := mapperOrm()
	clone.DB = orm.DB
	clone.sql = sql
	clone.args = args
	return clone
}

func (this *Orm) setDB(db *gorm.DB) {
	this.DB = db
}

func (this *Orm) Query() *gorm.DB {
	sqlReplace(this.NewScope(this))
	return this.DB.Raw(this.sql, this.args...)
}

func (this *Orm) Exec() *gorm.DB {
	sqlReplace(this.NewScope(this))
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

type Pager struct {
	Total,
	CurrentPage,
	PerPage,
	LastPage uint64
	Items interface{}
}

func NewPager(perPage, page uint64) *Pager {
	return &Pager{CurrentPage: page, PerPage: perPage}
}

func (this *Orm) Paginate(pager *Pager, items interface{}) {
	if pager.CurrentPage <= 0 {
		panic("Current Page: Can't be less than or equal to 0")
	}
	if pager.PerPage <= 0 {
		panic("Per Page: Can't be less than or equal to 0")
	}
	this.DB.
		Count(&pager.Total).
		Limit(pager.PerPage).
		Offset((pager.CurrentPage - 1) * pager.PerPage).
		Find(items)
	if this.DB.Error != nil {
		panic(this.DB.Error)
	}
	pager.LastPage = uint64(math.Ceil(float64(pager.Total) / float64(pager.PerPage)))
}

func (this *Orm) Table(name string) *Orm {
	this.DB = this.DB.Table(name)
	return this
}

func (this *Orm) Where(query interface{}, args ...interface{}) *Orm {
	this.DB = this.DB.Where(query, args...)
	return this
}

func (this *Orm) Or(query interface{}, args ...interface{}) *Orm {
	this.DB = this.DB.Or(query, args...)
	return this
}

func (this *Orm) Not(query interface{}, args ...interface{}) *Orm {
	this.DB = this.DB.Not(query, args...)
	return this
}

func (this *Orm) Limit(limit interface{}) *Orm {
	this.DB = this.DB.Not(limit)
	return this
}

func (this *Orm) Offset(offset interface{}) *Orm {
	this.DB = this.DB.Not(offset)
	return this
}

func (this *Orm) Order(value interface{}, reorder ...bool) *Orm {
	this.DB = this.DB.Order(value, reorder...)
	return this
}

func (this *Orm) Select(query interface{}, args ...interface{}) *Orm {
	this.DB = this.DB.Select(query, args...)
	return this
}

func (this *Orm) Omit(columns ...string) *Orm {
	this.DB = this.DB.Omit(columns...)
	return this
}

func (this *Orm) Group(query string) *Orm {
	this.DB = this.DB.Group(query)
	return this
}

func (this *Orm) Having(query interface{}, values ...interface{}) *Orm {
	this.DB = this.DB.Having(query, values...)
	return this
}

func (this *Orm) Joins(query string, args ...interface{}) *Orm {
	this.DB = this.DB.Joins(query, args...)
	return this
}

func (this *Orm) First(out interface{}, where ...interface{}) *Orm {
	this.DB = this.DB.First(out, where...)
	return this
}

func (this *Orm) Take(out interface{}, where ...interface{}) *Orm {
	this.DB = this.DB.Take(out, where...)
	return this
}

func (this *Orm) Last(out interface{}, where ...interface{}) *Orm {
	this.DB = this.DB.Last(out, where...)
	return this
}

func (this *Orm) Find(out interface{}, where ...interface{}) *Orm {
	this.DB = this.DB.Find(out, where...)
	return this
}

func (this *Orm) Pluck(column string, value interface{}) *Orm {
	this.DB = this.DB.Pluck(column, value)
	return this
}

func (this *Orm) Count(value interface{}) *Orm {
	this.DB = this.DB.Count(value)
	return this
}
