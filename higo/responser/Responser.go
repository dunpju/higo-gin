package responser

import "github.com/dengpju/higo-gin/higo"

func SuccessJson(message, code interface{}, data interface{}) {
	ctx := higo.Request.Context()
	higo.Responser(ctx).SuccessJson(message.(string), code.(int), data)
}

func ErrorJson(message, code interface{}, data interface{}) {
	ctx := higo.Request.Context()
	higo.Responser(ctx).ErrorJson(message.(string), code.(int), data)
}
