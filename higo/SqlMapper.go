package higo

import (
	"github.com/dengpju/higo-throw/exception"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type SqlMapper struct {
	Orm  *Orm
	Sql  string
	Args []interface{}
	orms []*Orm
}

func Mapper(sql string, args []interface{}, err error) *SqlMapper {
	if err != nil {
		panic(err.Error())
	}
	sqlMapper := &SqlMapper{}
	sqlMapper.Orm = newOrm()
	sqlMapper.Sql = sql
	sqlMapper.Args = args
	sqlMapper.Orm.sql = sql
	sqlMapper.Orm.args = args
	return sqlMapper
}

func (this *SqlMapper) Query() *gorm.DB {
	return this.Orm.DB.Raw(this.Sql, this.Args)
}

func (this *SqlMapper) Exec() *gorm.DB {
	return this.Orm.DB.Exec(this.Sql, this.Args)
}

func (this *SqlMapper) setDB(db *gorm.DB) {
	this.Orm.DB = db
}

func (this *SqlMapper) apply(tx *gorm.DB) {
	for _, o := range this.orms {
		o.setDB(tx)
	}
}

func (this *SqlMapper) Transaction(fn func() error) {
	err := this.Orm.DB.Transaction(func(tx *gorm.DB) (err error) {
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
		exception.Throw(exception.Message(err.(*mysql.MySQLError).Message),
			exception.Code(int(err.(*mysql.MySQLError).Number)),
			exception.Data(nil))
	}
}

func Begin(orms ...*Orm) *SqlMapper {
	sqlMapper := &SqlMapper{}
	sqlMapper.Orm = newOrm()
	sqlMapper.orms = orms
	return sqlMapper
}
