package sql

import (
	"github.com/Masterminds/squirrel"
)

var StatementBuilder = squirrel.StatementBuilder

func Select(columns ...string) squirrel.SelectBuilder {
	return squirrel.Select(columns...)
}

func Insert(into string) squirrel.InsertBuilder {
	return squirrel.Insert(into)
}

func Replace(into string) squirrel.InsertBuilder {
	return squirrel.Replace(into)
}

func Update(table string) squirrel.UpdateBuilder {
	return squirrel.Update(table)
}

func Delete(from string) squirrel.DeleteBuilder {
	return squirrel.Delete(from)
}

func Case(what ...interface{}) squirrel.CaseBuilder {
	return squirrel.Case(what)
}
