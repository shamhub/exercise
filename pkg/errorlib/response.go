package errorlib

import "fmt"

type ResponseError struct {
	code   int
	reason string
}

func NewResponseError(code int, reason string) *ResponseError {
	return &ResponseError{
		code:   code,
		reason: reason,
	}
}

func (err *ResponseError) Error() string {
	code := fmt.Sprintf("code: %d", err.code)
	reason := fmt.Sprintf("reason: %s", err.reason)
	return code + "-" + reason
}

func (err *ResponseError) GetStatusCode() int {
	return err.code
}

func (err *ResponseError) ProvideReason() string {
	return err.reason
}
