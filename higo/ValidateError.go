package higo

import "gitee.com/dengpju/higo-code/code"

type ValidateError struct {
	error code.ICode
}

func NewValidateError(error code.ICode) *ValidateError {
	return &ValidateError{error: error}
}

func (this ValidateError) Error() string {
	return this.error.Message()
}

func (this ValidateError) Get() code.ICode {
	return this.error
}
