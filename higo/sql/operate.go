package sql

import (
	"github.com/dengpju/higo-utils/utils"
	"strings"
)

func convert(values interface{}) []string {
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
	return column + " IN(" + strings.Join(convert(values), ",") + ")"
}

func NotIn(column string, values interface{}) string {
	return column + " NOT IN(" + strings.Join(convert(values), ",") + ")"
}

func IsNull(column string) string {
	return "isnull(`" + column + "`)"
}

func IF(expr1, expr2, expr3 string) string {
	return "IF(" + expr1 + "," + expr2 + "," + expr3 + ")"
}
