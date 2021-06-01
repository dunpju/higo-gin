package higo

import (
	"github.com/jinzhu/gorm"
)

type ISqlMapper interface {
	Sql() string
	Args() []interface{}
	Orm() interface{}
}

func Mapper(sql string, args []interface{}, err error) *GormSqlMapperImpl {
	if err != nil {
		panic(err.Error())
	}
	return NewGormSqlMapper(sql, args)
}

type SqlMapperAbstract struct {
	orm  *Orm
	sql  string
	args []interface{}
}

func NewSqlMapperAbstract(orm *Orm, sql string, args []interface{}) *SqlMapperAbstract {
	return &SqlMapperAbstract{orm: orm, sql: sql, args: args}
}

func (this *SqlMapperAbstract) Sql() string {
	return this.sql
}

func (this *SqlMapperAbstract) Args() []interface{} {
	return this.args
}

func (this *SqlMapperAbstract) Orm() interface{} {
	return this.orm
}

type GormSqlMapperImpl struct {
	*SqlMapperAbstract
}

func NewGormSqlMapper(sql string, args []interface{}) *GormSqlMapperImpl {
	return &GormSqlMapperImpl{NewSqlMapperAbstract(NewOrm(), sql, args)}
}

func (this *GormSqlMapperImpl) Query() *gorm.DB {
	return this.orm.Raw(this.sql, this.Args)
}

func (this *GormSqlMapperImpl) Exec() *gorm.DB {
	return this.orm.Exec(this.sql, this.Args)
}
