package sql

import (
	"github.com/dengpju/higo-utils/utils"
	"strings"
)

func WhereIn(column string, values []interface{}) string {
	conValues := make([]string, 0)
	for _, value := range values {
		if v, ok := value.(int); ok {
			conValues = append(conValues, utils.IntString(v))
		} else if v, ok := value.(int64); ok {
			conValues = append(conValues, utils.Int64String(v))
		} else if v, ok := value.(float32); ok {
			conValues = append(conValues, utils.FloatString(v))
		} else if v, ok := value.(float64); ok {
			conValues = append(conValues, utils.Float64String(v))
		} else {
			conValues = append(conValues, value.(string))
		}
	}
	return column + " IN(" + strings.Join(conValues, ",") + ")"
}

func WhereNotIn(column string, values []interface{}) string {
	conValues := make([]string, 0)
	for _, value := range values {
		if v, ok := value.(int); ok {
			conValues = append(conValues, utils.IntString(v))
		} else if v, ok := value.(int64); ok {
			conValues = append(conValues, utils.Int64String(v))
		} else if v, ok := value.(float32); ok {
			conValues = append(conValues, utils.FloatString(v))
		} else if v, ok := value.(float64); ok {
			conValues = append(conValues, utils.Float64String(v))
		} else {
			conValues = append(conValues, value.(string))
		}
	}
	return column + " NOT IN(" + strings.Join(conValues, ",") + ")"
}
