package sql

import (
	"fmt"
	"strings"
)

func Field(fields ...string) string {
	if len(fields) == 0 {
		panic(fmt.Errorf("Fields Can Not Be Empty"))
	}
	tmpFields := make([]string, 0)
	for _, f := range fields {
		f = strings.ReplaceAll(f, "`", "")
		fs := strings.Split(f, ".")
		if len(fs) >= 2 {
			tmpFields = append(tmpFields, "`"+fs[0]+"`.`"+fs[1]+"`")
		} else {
			tmpFields = append(tmpFields, "`"+fs[0]+"`")
		}
	}
	return strings.Join(tmpFields, ",")
}
