package higo

type ValidateError struct {
	error string
}

func NewValidateError(error string) *ValidateError {
	return &ValidateError{error: error}
}

func (this ValidateError) Error() string {
	return this.error
}
