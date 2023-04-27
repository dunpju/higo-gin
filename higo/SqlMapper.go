package higo

import (
	"github.com/jinzhu/gorm"
)

type SqlMapper struct {
	Orm  *Orm
	orms []*Orm
}

func NewSqlMapper(orm *Orm, orms []*Orm) *SqlMapper {
	return &SqlMapper{Orm: orm, orms: orms}
}

func Mapper(sql string, args []interface{}, err error) *SqlMapper {
	if err != nil {
		panic(err.Error())
	}
	sqlMapper := &SqlMapper{}
	sqlMapper.Orm = newOrm()
	sqlMapper.Orm.sql = sql
	sqlMapper.Orm.args = args
	return sqlMapper
}

func (this *SqlMapper) Query() *gorm.DB {
	return this.Orm.DB.Raw(this.Orm.sql, this.Orm.args)
}

func (this *SqlMapper) Exec() *gorm.DB {
	return this.Orm.DB.Exec(this.Orm.sql, this.Orm.args)
}

func (this *SqlMapper) setDB(db *gorm.DB) {
	this.Orm.DB = db
}

func (this *SqlMapper) Transaction(fn func() error) {
	this.Orm.Begin(this.orms...).Transaction(fn)
}

func Begin(orms ...*Orm) *SqlMapper {
	return NewSqlMapper(newOrm(), orms)
}
