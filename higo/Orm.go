package higo

import (
	"database/sql"
	"fmt"
	"github.com/dengpju/higo-config/config"
	higosql "github.com/dengpju/higo-gin/higo/sql"
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
	"time"
)

var (
	orm                     *Orm
	onceGorm                sync.Once
	dbConfigOnce            sync.Once
	confDefault             *config.Configure
	dbConfig                *Dbconfig
	logMode                 bool
	maxIdle                 = 5
	maxOpen                 = 5
	maxLifetime             = 20
	registerCallbackCounter int
)

func GetDbConfig() *Dbconfig {
	return dbConfig
}

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

// gorm实例
func newGorm() *gorm.DB {
	dbConfigOnce.Do(func() {
		confDefault = config.Db("DB.Default").(*config.Configure)
		dbConfig = &Dbconfig{Username: confDefault.Get("Username").(string),
			Password: confDefault.Get("Password").(string),
			Host:     confDefault.Get("Host").(string),
			Port:     confDefault.Get("Port").(string),
			Database: confDefault.Get("Database").(string),
			Charset:  confDefault.Get("Charset").(string),
			Driver:   confDefault.Get("Driver").(string),
			Prefix:   confDefault.Get("Prefix").(string),
		}
		logMode = confDefault.Get("LogMode").(bool)
		maxIdle = confDefault.Get("MaxIdle").(int)
		maxOpen = confDefault.Get("MaxOpen").(int)
		maxLifetime = confDefault.Get("MaxLifetime").(int)
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
	db.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
	if registerCallbackCounter == 1 {
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
		logger.Logrus.Infoln(fmt.Sprintf("DB %s:%s Connection success!", dbConfig.Host,
			dbConfig.Port))
	}
	return db
}

func sqlReplace(scope *gorm.Scope) {
	sq := scope.SQL
	s := reflect.ValueOf(scope.SQLVars)
	if orm, ok := scope.Value.(*Orm); ok {
		sq = orm.sql
		s = reflect.ValueOf(orm.args)
	}
	for i := 0; i < s.Len(); i++ {
		sq = strings.Replace(sq, "?", "'%v'", 1)
		sq = fmt.Sprintf(sq, s.Index(i))
	}
	if logMode {
		logger.Logrus.Debugln(sq)
	}
}

func NewOrm() *Orm {
	return newOrm()
}

// 单例
func SingleOrm() *Orm {
	onceGorm.Do(func() {
		orm = newOrm()
	})
	return orm
}

func newOrm() *Orm {
	registerCallbackCounter++
	return &Orm{DB: newGorm(), orms: make([]*Orm, 0),
		builder: newBuilder(),
	}
}

type Orm struct {
	*gorm.DB
	sql          string
	args         []interface{}
	orms         []*Orm
	table        string
	statement    *higosql.Statement
	builder      *Builder
	result       sql.Result
	lastInsertId int64
}

func (this *Orm) Args() []interface{} {
	return this.args
}

func (this *Orm) Sql() string {
	return this.sql
}

func (this *Orm) Mapper(sql string, args []interface{}, err error) *Orm {
	if err != nil {
		panic(err.Error())
	}
	this.sql = sql
	this.args = args
	return this
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

func (this *Orm) InsertGetId() *Orm {
	sqlReplace(this.NewScope(this))
	this.result, this.DB.Error = this.DB.CommonDB().Exec(this.sql, this.args...)
	return this
}

func (this *Orm) LastInsertId() int64 {
	id, err := this.result.LastInsertId()
	if nil != err {
		panic(err)
	}
	return id
}

func (this *Orm) Result() sql.Result {
	return this.result
}

func (this *Orm) Begin(orms ...*Orm) *Orm {
	this.orms = orms
	return this
}

func (this *Orm) apply(tx *gorm.DB) {
	for _, s := range this.orms {
		s.setDB(tx)
	}
}

// 事务
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

func (this *Orm) Table(name string) *Orm {
	this.DB = this.DB.New().Table(name)
	this.table = name
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

func (this *Orm) Builder() *higosql.Statement {
	this.statement = higosql.Query()
	return this.statement
}

func (this *Orm) Insert(name string) *higosql.Statement {
	this.table = name
	this.statement = higosql.Insert(name)
	return this.statement
}

func (this *Orm) Update(name string) *higosql.Statement {
	this.table = name
	this.statement = higosql.Update(name)
	return this.statement
}

func (this *Orm) Delete(name string) *higosql.Statement {
	this.table = name
	this.statement = higosql.Delete(name)
	return this.statement
}

func (this *Orm) ToSql() (string, []interface{}, error) {
	return this.statement.ToSql()
}

func (this *Orm) setBuilder(sql string, args []interface{}, err error) {
	this.builder.sql, this.builder.args, this.builder.err = sql, args, err
}

func (this *Orm) GetBuilder() (string, []interface{}, error) {
	return this.builder.sql, this.builder.args, this.builder.err
}

func (this *Orm) Build() {
	this.setBuilder(this.ToSql())
}

func (this *Orm) CheckError() {
	if this.Error != nil {
		panic(this.Error)
	}
}

func (this *Orm) Paginate(pager *Pager) *gorm.DB {
	if pager.CurrentPage <= 0 {
		panic("Current Page: Can't be less than or equal to 0")
	}
	if pager.PerPage <= 0 {
		panic("Per Page: Can't be less than or equal to 0")
	}
	this.DB.Count(&pager.Total)
	if this.DB.Error != nil {
		panic(this.DB.Error)
	}
	pager.LastPage = uint64(math.Ceil(float64(pager.Total) / float64(pager.PerPage)))
	return this.DB.Limit(pager.PerPage).
		Offset((pager.CurrentPage - 1) * pager.PerPage)
}

func NewPager(perPage, page uint64) *Pager {
	return &Pager{CurrentPage: page, PerPage: perPage}
}

type Pager struct {
	Total       uint64      `json:"total"`
	CurrentPage uint64      `json:"current_page"`
	PerPage     uint64      `json:"per_page"`
	LastPage    uint64      `json:"last_page"`
	Items       interface{} `json:"items"`
}

func newBuilder() *Builder {
	return &Builder{}
}

type Builder struct {
	sql  string
	args []interface{}
	err  error
}
