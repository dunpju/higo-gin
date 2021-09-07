package sql

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

type operationState int

const (
	opSelect operationState = iota + 1
	opInsert
	opUpdate
	opDelete
)

type setClause struct {
	column string
	value  interface{}
}

type setWhere struct {
	pred  string
	value interface{}
}

type Statement struct {
	Builder        interface{}
	currentOpState operationState
	setClauses     []setClause
	setWheres      []setWhere
	table          string
}

func statement(builder interface{}, opState operationState, table string) *Statement {
	return &Statement{Builder: builder, currentOpState: opState, table: table}
}

func (this *Statement) Set(column string, value interface{}) *Statement {
	this.setClauses = append(this.setClauses, setClause{column: column, value: value})
	return this
}

func (this *Statement) Where(pred string, value interface{}) *Statement {
	this.setWheres = append(this.setWheres, setWhere{pred: pred, value: value})
	return this
}

func (this *Statement) SelectBuilder() squirrel.SelectBuilder {
	return this.Builder.(squirrel.SelectBuilder)
}

func (this *Statement) InsertBuilder() squirrel.InsertBuilder {
	return this.Builder.(squirrel.InsertBuilder)
}

func (this *Statement) UpdateBuilder() squirrel.UpdateBuilder {
	return this.Builder.(squirrel.UpdateBuilder)
}

func (this *Statement) DeleteBuilder() squirrel.DeleteBuilder {
	return this.Builder.(squirrel.DeleteBuilder)
}

func (this *Statement) ToSql() (string, []interface{}, error) {
	var (
		columns []string
		values  []interface{}
	)
	if opSelect == this.currentOpState {
		return this.Builder.(squirrel.SelectBuilder).ToSql()
	} else if opInsert == this.currentOpState {
		for _, clause := range this.setClauses {
			columns = append(columns, clause.column)
			values = append(values, clause.value)
		}
		return this.Builder.(squirrel.InsertBuilder).Columns(columns...).Values(values...).ToSql()
	} else if opUpdate == this.currentOpState {
		setMap := make(map[string]interface{}, 0)
		for _, clause := range this.setClauses {
			setMap[clause.column] = clause.value
		}
		whereMap := make(map[string]interface{}, 0)
		for _, where := range this.setWheres {
			whereMap[where.pred] = where.value
		}
		return this.Builder.(squirrel.UpdateBuilder).SetMap(setMap).Where(whereMap).ToSql()
	} else if opDelete == this.currentOpState {
		whereMap := make(map[string]interface{}, 0)
		for _, where := range this.setWheres {
			whereMap[where.pred] = where.value
		}
		return this.Builder.(squirrel.DeleteBuilder).Where(whereMap).ToSql()
	}
	return "", nil, fmt.Errorf("An unsupported operation")
}

func Query() *Statement {
	return statement(squirrel.StatementBuilder, opSelect, "")
}

func (this *Statement) Select(columns ...string) squirrel.SelectBuilder {
	return statement(Select(columns...), opSelect, "").SelectBuilder()
}

func Select(columns ...string) squirrel.SelectBuilder {
	return Query().Select(columns...)
}

func Insert(into string) *Statement {
	return statement(squirrel.Insert(into), opInsert, into)
}

func Update(table string) *Statement {
	return statement(squirrel.Update(table), opUpdate, table)
}

func Delete(from string) *Statement {
	return statement(squirrel.Delete(from), opDelete, from)
}
