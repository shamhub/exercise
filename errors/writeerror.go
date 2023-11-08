package errors

import "net/http"

type WriteError struct {
	Code   string
	Reason string
}

func NewWriteError(code, reason string) *WriteError {
	return &WriteError{
		Code:   code,
		Reason: reason,
	}
}

func (error *WriteError) Error() string {
	return error.Code + " - " + error.Reason
}

type MultiErrors = []*WriteError

func NewMultipleProcessError(mErrs MultiErrors) error {
	errs := make([]error, len(mErrs))
	for i, e := range mErrs {
		errs[i] = NewCustomError(e.Code, e.Reason, http.StatusBadRequest)
	}
	resp := MultipleErrors{StatusCode: http.StatusBadRequest, Errors: errs}
	return resp
}
