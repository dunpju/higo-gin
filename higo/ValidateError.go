package higo

import "gitee.com/dengpju/higo-code/code"

type ValidateError struct {
	error *code.Code
}

func NewValidateError(error *code.Code) *ValidateError {
	return &ValidateError{error: error}
}

func (this ValidateError) Error() string {
	return this.error.String()
}

func (this ValidateError) Get() *code.Code {
	return this.error
}
