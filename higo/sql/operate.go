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

type ConditionPrepare func() string

func Perd(cps ...ConditionPrepare) string {
	if len(cps) == 0 {
		panic("Raw Condition Can Not Be Empty")
	}
	condSlice := make([]string, 0)
	for _, cond := range cps {
		condSlice = append(condSlice, cond())
	}
	return "(" + strings.Join(condSlice, " AND ") + ")"
}

func AND(cps ...ConditionPrepare) ConditionPrepare {
	return func() string {
		if len(cps) == 0 {
			panic("AND Condition Can Not Be Empty")
		}
		condSlice := make([]string, 0)
		for _, cond := range cps {
			condSlice = append(condSlice, cond())
		}
		return "(" + strings.Join(condSlice, " AND ") + ")"
	}
}

func OR(cps ...ConditionPrepare) ConditionPrepare {
	return func() string {
		if len(cps) == 0 {
			panic("OR Condition Can Not Be Empty")
		}
		condSlice := make([]string, 0)
		for _, cond := range cps {
			condSlice = append(condSlice, cond())
		}
		return "(" + strings.Join(condSlice, " OR ") + ")"
	}
}

// where Condition
func Condition(column, operator string, value interface{}) ConditionPrepare {
	return func() string {
		column = strings.ReplaceAll(column, "`", "")
		columns := strings.Split(column, ".")
		if len(columns) >= 2 {
			return "(`" + columns[0] + "`.`" + columns[1] + "` " + operator + `'` + convertString(value) + `')`
		} else {
			return "(`" + column + "` " + operator + `'` + convertString(value) + `')`
		}
	}
}

func BETWEEN(column string, value1, value2 interface{}) ConditionPrepare {
	return func() string {
		return column + ` BETWEEN ` + convertString(value1) + ` AND ` + convertString(value2)
	}
}

func IN(column string, values interface{}) ConditionPrepare {
	return func() string {
		return column + " IN(" + strings.Join(convertSliceString(values), ",") + ")"
	}
}

func NOTIN(column string, values interface{}) ConditionPrepare {
	return func() string {
		return column + " NOT IN(" + strings.Join(convertSliceString(values), ",") + ")"
	}
}

func ISNULL(column string) ConditionPrepare {
	return func() string {
		column = strings.ReplaceAll(column, "`", "")
		return "(`" + column + "`" + ` IS NULL)`
	}
}

func Raw(query string) ConditionPrepare {
	return func() string {
		return query
	}
}
