package sql

import "strings"

type Join func() string

func Joins(joins ...Join) string {
	if len(joins) == 0 {
		panic("Joins Can Not Be Empty")
	}
	joinSlice := make([]string, 0)
	for _, j := range joins {
		joinSlice = append(joinSlice, j())
	}
	return strings.Join(joinSlice, " ")
}

func LeftJoin(table, first, operator, second string) Join {
	return func() string {
		return "LEFT JOIN " + table + " ON " + first + operator + second
	}
}

func RightJoin(table, first, operator, second string) Join {
	return func() string {
		return "RIGHT JOIN " + table + " ON " + first + operator + second
	}
}

func InnerJoin(table, first, operator, second string) Join {
	return func() string {
		return "INNER JOIN " + table + " ON " + first + operator + second
	}
}
