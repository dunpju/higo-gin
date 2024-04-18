package higo

import (
	"encoding/json"
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

type IResult interface {
	SetCode(code int)
	SetMessage(msg string)
	SetData(data interface{})
}

type JsonResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (j *JsonResult) SetCode(code int) {
	j.Code = code
}

func (j *JsonResult) SetMessage(msg string) {
	j.Message = msg
}

func (j *JsonResult) SetData(data interface{}) {
	j.Data = data
}

func NewJsonResult(code int, message string, data interface{}) IResult {
	return &JsonResult{Code: code, Message: message, Data: data}
}

type ResultHandler func(code int, message string, data interface{}) IResult

var (
	resultPool *sync.Pool
	NewResult  ResultHandler = NewJsonResult
)

func init() {
	resultPool = &sync.Pool{
		New: func() interface{} {
			return NewResult(0, "", nil)
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
		result := resultPool.Get().(IResult)
		defer resultPool.Put(result)
		result.SetMessage(message)
		result.SetCode(code)
		result.SetData(data)
		return func(output Output) {
			output(ctx, result)
		}
	}
}

func ResponserTest() ResultFunc {
	return func(message string, code int, data interface{}) func(output Output) {
		result := resultPool.Get().(IResult)
		defer resultPool.Put(result)
		result.SetMessage(message)
		result.SetCode(code)
		result.SetData(data)
		fmt.Println(message, code, data)
		fmt.Println(result)
		marshal, err := json.Marshal(result)
		if err != nil {
			return nil
		}
		fmt.Println(string(marshal))
		return func(output Output) {}
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

// Receiver 结果接收者
func Receiver(values ...interface{}) *ErrorResult {
	return Result(values...)
}
