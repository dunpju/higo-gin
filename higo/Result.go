package higo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type ErrorResult struct {
	data interface{}
	err  error
}

func (this *ErrorResult) Unwrap() interface{} {
	if this.err != nil {
		checkErrors(this.err)
		panic(this.err)
	}
	return this.data
}

func (this *ErrorResult) Error() error {
	return this.err
}

func Result(values ...interface{}) *ErrorResult {
	if len(values) == 1 {
		if values[0] == nil {
			return &ErrorResult{nil, nil}
		}
		if e, ok := values[0].(error); ok {
			return &ErrorResult{nil, e}
		}
	} else if len(values) == 2 {
		if values[1] == nil {
			return &ErrorResult{values[0], nil}
		}
		if e, ok := values[1].(error); ok {
			return &ErrorResult{values[0], e}
		}
	}
	return &ErrorResult{nil, fmt.Errorf("error result format")}
}

func checkErrors(errs error) {
	if errs, ok := errs.(ValidateError); ok {
		panic(errs)
	}
}

type JsonResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewJsonResult(code int, message string, data interface{}) *JsonResult {
	return &JsonResult{Code: code, Message: message, Data: data}
}

var resultPool *sync.Pool

func init() {
	resultPool = &sync.Pool{
		New: func() interface{} {
			return NewJsonResult(0, "", nil)
		},
	}
}

type ResultFunc func(message string, code int, data interface{}) func(output Output)
type Output func(ctx *gin.Context, v interface{})

func (this ResultFunc) SuccessJson(message string, code int, data interface{}) {
	this(message, code, data)(OK)
}

func (this ResultFunc) ErrorJson(message string, code int, data interface{}) {
	this(message, code, data)(Error)
}

func Responser(ctx *gin.Context) ResultFunc {
	return func(message string, code int, data interface{}) func(output Output) {
		r := resultPool.Get().(*JsonResult)
		defer resultPool.Put(r)
		r.Message = message
		r.Code = code
		r.Data = data
		return func(output Output) {
			output(ctx, r)
		}
	}
}

func OK(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, v)
	panic(nil)
}

func Error(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusBadRequest, v)
	panic(nil)
}

//结果接收者
func Receiver(values ...interface{}) *ErrorResult {
	return Result(values...)
}
