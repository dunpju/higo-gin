package higo

import (
	"github.com/jinzhu/gorm"
)

type SqlMapper struct {
	Orm  *Orm
	Sql  string
	Args []interface{}
}

func Mapper(sql string, args []interface{}, err error) *SqlMapper {
	if err != nil {
		panic(err.Error())
	}
	sqlMapper := &SqlMapper{}
	sqlMapper.Orm = newOrm()
	sqlMapper.Sql = sql
	sqlMapper.Args = args
	return sqlMapper
}

func (this *SqlMapper) Query() *gorm.DB {
	return this.Orm.Raw(this.Sql, this.Args)
}

func (this *SqlMapper) Exec() *gorm.DB {
	return this.Orm.Exec(this.Sql, this.Args)
}

func (this *SqlMapper) setDB(db *gorm.DB) {
	this.Orm.DB = db
}

func (this *SqlMapper) apply(orms ...*Orm) {
	for _, orm := range orms {
		orm.setDB(this.Orm.DB)
	}
}

func (this *SqlMapper) Transaction(fn func() error) {
	err := this.Orm.DB.Transaction(func(tx *gorm.DB) error {
		return fn()
	})
	panic(err)
}

func Begin(orms ...*Orm) *SqlMapper {
	sqlMapper := &SqlMapper{}
	sqlMapper.Orm = newOrm()
	sqlMapper.apply(orms...)
	return sqlMapper
}
