package higo

import (
	"fmt"
	"gitee.com/dengpju/higo-parameter/parameter"
	"github.com/dengpju/higo-throw/exception"
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
