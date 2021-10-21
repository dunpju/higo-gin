package sql

func IsNull(column string) string {
	return "isnull(`" + column + "`)"
}

func If(expr1, expr2, expr3 string) string {
	return "if(" + expr1 + "," + expr2 + "," + expr3 + ")"
}
