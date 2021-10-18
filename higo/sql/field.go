package sql

import "strings"

func Field(fields ...string) string {
	if len(fields) > 0 {
		return strings.Join(fields, ",")
	}
	return "*"
}
