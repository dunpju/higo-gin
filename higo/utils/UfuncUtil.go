package utils

// 三目运算
func If(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}
