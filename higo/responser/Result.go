package responser

import (
	"fmt"
	"gitee.com/dengpju/higo-parameter/parameter"
	"github.com/dengpju/higo-throw/exception"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type ErrorResult struct {
	data interface{}
	err  error
}

func (this *ErrorResult) Unwrap(parameters ...*parameter.Parameter) interface{} {
	if this.err != nil {
		if len(parameters) > 0 {
			tmp := make([]*parameter.Parameter, 0)
			for _, p := range parameters {
				if p.Name == exception.MESSAGE {
					tmp = append(tmp, p)
					tmp = append(tmp, exception.RealMessage(this.err))
				} else {
					tmp = append(tmp, p)
				}
			}
			exception.Throw(exception.Message(this.err))
		} else {
			exception.Throw(exception.Message(this.err), exception.RealMessage(this.err))
		}
	}
	return this.data
}

func Result(values ...interface{}) *ErrorResult {
	if len(values) == 1 {
		if values[0] == nil {
			return &ErrorResult{nil, nil}
		}
		if e, ok := values[0].(error); ok {
			return &ErrorResult{nil, e}
		}
	}
	if len(values) == 2 {
		if values[1] == nil {
			return &ErrorResult{values[0], nil}
		}
		if e, ok := values[1].(error); ok {
			return &ErrorResult{values[0], e}
		}
	}
	return &ErrorResult{nil, fmt.Errorf("error result format")}
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

func End(ctx *gin.Context) ResultFunc {
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
