package higo

// 异常接口
type IThrowable interface {
	Throw(message string, code int)
}