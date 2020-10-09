package higo

// 异常接口
type IThrowable interface {
	Exception(message interface{}, code int, data ...interface{})
}