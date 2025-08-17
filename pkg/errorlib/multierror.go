package errorlib

import "fmt"

type MultiErrors struct {
	statusCode int
	errors     []*ResponseError
}

func NewMultiError(code int, responseErrorList []*ResponseError) *MultiErrors {
	errList := make([]*ResponseError, len(responseErrorList))
	copy(errList, responseErrorList)
	return &MultiErrors{
		statusCode: code,
		errors:     errList,
	}
}

func (m *MultiErrors) Error() string {
	if m == nil {
		panic("multierror object is nil")
	}
	code := fmt.Sprintf("status code: %d, ", m.statusCode)
	var errorString string
	for index, e := range m.errors {
		str := fmt.Sprintf("error %d - %s", index, e.Error())
		errorString += str
	}
	return code + errorString
}

func (m *MultiErrors) ProvideReason() []string {
	reasonList := make([]string, len(m.errors))
	for index, e := range m.errors {
		str := fmt.Sprintf("error %d - %s", index, e.Error())
		reasonList[index] += str
	}
	return reasonList
}

func (m *MultiErrors) GetStatusCode() int {
	return m.statusCode
}
