package sql

import (
	"github.com/dengpju/higo-utils/utils"
	"strings"
)

func convertString(values interface{}) string {
	if value, ok := values.(int); ok {
		return utils.IntString(value)
	} else if value, ok := values.(int64); ok {
		return utils.Int64String(value)
	} else if value, ok := values.(float32); ok {
		return utils.FloatString(value)
	} else if value, ok := values.(float64); ok {
		utils.Float64String(value)
	} else if value, ok := values.(string); ok {
		return value
	} else {
		panic("Unsupported types")
	}
	return ""
}

func convertSliceString(values interface{}) []string {
	conValues := make([]string, 0)
	if value, ok := values.([]int); ok {
		for _, v := range value {
			conValues = append(conValues, utils.IntString(v))
		}
	} else if value, ok := values.([]int64); ok {
		for _, v := range value {
			conValues = append(conValues, utils.Int64String(v))
		}
	} else if value, ok := values.([]float32); ok {
		for _, v := range value {
			conValues = append(conValues, utils.FloatString(v))
		}
	} else if value, ok := values.([]float64); ok {
		for _, v := range value {
			conValues = append(conValues, utils.Float64String(v))
		}
	} else if value, ok := values.([]string); ok {
		conValues = value
	} else {
		panic("Unsupported types")
	}
	return conValues
}

func IN(column string, values interface{}) string {
	return column + " IN(" + strings.Join(convertSliceString(values), ",") + ")"
}

func NotIn(column string, values interface{}) string {
	return column + " NOT IN(" + strings.Join(convertSliceString(values), ",") + ")"
}

func IsNull(column string) string {
	return "isnull(`" + column + "`)"
}

func IF(expr1, expr2, expr3 string) string {
	return "IF(" + expr1 + "," + expr2 + "," + expr3 + ")"
}

type Prep func() string

func Perd(conds ...Prep) string {
	if len(conds) == 0 {
		panic("Raw Condition Can Not Be Empty")
	}
	condSlice := make([]string, 0)
	for _, cond := range conds {
		condSlice = append(condSlice, cond())
	}
	return "(" + strings.Join(condSlice, " AND ") + ")"
}

func Raw(query string) Prep {
	return func() string {
		return query
	}
}

func AND(conds ...Prep) Prep {
	return func() string {
		if len(conds) == 0 {
			panic("AND Condition Can Not Be Empty")
		}
		condSlice := make([]string, 0)
		for _, cond := range conds {
			condSlice = append(condSlice, cond())
		}
		return "(" + strings.Join(condSlice, " AND ") + ")"
	}
}

func OR(conds ...Prep) Prep {
	return func() string {
		if len(conds) == 0 {
			panic("OR Condition Can Not Be Empty")
		}
		condSlice := make([]string, 0)
		for _, cond := range conds {
			condSlice = append(condSlice, cond())
		}
		return "(" + strings.Join(condSlice, " OR ") + ")"
	}
}

func Cond(column, operator string, value interface{}) Prep {
	return func() string {
		columns := strings.Split(column, ".")
		if len(columns) >= 2 {
			return "(`" + columns[0] + "`.`" + columns[1] + "` " + operator + `'` + convertString(value) + `')`
		} else {
			return "(`" + column + "` " + operator + `'` + convertString(value) + `')`
		}
	}
}
