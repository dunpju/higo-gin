package mapper

import "github.com/jinzhu/gorm"

type SqlMapper struct {
	Sql  string
	Args []interface{}
}

func NewSqlMapper(sql string, args []interface{}) *SqlMapper {
	return &SqlMapper{Sql: sql, Args: args}
}

func (this *SqlMapper) Query() *gorm.DB {
 return dbs.Orm.Raw(this.Sql, this.Args)
}

func (this *SqlMapper) Exec() *gorm.DB {
	return dbs.Orm.Exec(this.Sql, this.Args)
}

func Mapper(sql string, args []interface{}, err error) *SqlMapper {
	if err != nil {
		panic(err.Error())
	}
	return NewSqlMapper(sql, args)
}
