package responser

import "github.com/dengpju/higo-gin/higo"

func SuccessJson(message string, code int, data interface{}) {
	ctx := higo.Request.Context()
	higo.Responser(ctx).SuccessJson(message, code, data)
}

func ErrorJson(message string, code int, data interface{}) {
	ctx := higo.Request.Context()
	higo.Responser(ctx).ErrorJson(message, code, data)
}
