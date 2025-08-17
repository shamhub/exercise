package httpservice

import (
	"encoding/json"
	"time"

	"github.com/shamhub/exercise/pkg/errorlib"
)

const (
	dateFormat = "2006-01-02T15:04:05.000Z"
)

type customError struct {
	StatusCode   int       `json:"statuscode" xml:"statuscode"`
	Reason       []string  `json:"reason" xml:"reason"`
	ResourcePath string    `json:"resourcepath" xml:"resourcepath"`
	RootCause    RootCause `json:"rootCauses,omitempty" xml:"root_causes,omitempty"`
	Title        string    `json:"title,omitempty" xml:"title,omitempty"`
	DateTime     DateTime  `json:"datetime" xml:"datetime"`
}

type RootCause map[string]string

type DateTime struct {
	Value    string `json:"value" xml:"value"`
	TimeZone string `json:"timezone" xml:"timezone"`
}

func (k *customError) Error() string {
	errorBytes, err := json.Marshal(k)
	if err != nil {
		panic(err)
	}
	return string(errorBytes)
}

func newCustomErrorForSingleErrorResponse(resourcePath string, err *errorlib.ResponseError) *customError {
	currentTime := time.Now().UTC()
	currentTimeFormatted := currentTime.Format(dateFormat)
	return &customError{
		StatusCode:   err.GetStatusCode(),
		Reason:       []string{err.ProvideReason()},
		ResourcePath: resourcePath,
		RootCause:    responseErrorCause(err),
		Title:        errorResponseTitle(err),
		DateTime: DateTime{
			Value:    currentTimeFormatted,
			TimeZone: "UTC",
		},
	}
}

func newCustomErrorForMultiErrorResponse(resourcePath string, err *errorlib.MultiErrors) *customError {
	currentTime := time.Now().UTC()
	currentTimeFormatted := currentTime.Format(dateFormat)
	return &customError{
		StatusCode:   err.GetStatusCode(),
		Reason:       err.ProvideReason(),
		ResourcePath: resourcePath,
		RootCause:    multiErrorCause(err),
		Title:        multiErrorResponseTitle(err),
		DateTime: DateTime{
			Value:    currentTimeFormatted,
			TimeZone: "UTC",
		},
	}
}

func (c *customError) GetStatusCode() int {
	return c.StatusCode
}

func responseErrorCause(err *errorlib.ResponseError) RootCause {
	// Interpret the error and provide possible cause to fix it
	return RootCause{}
}

func errorResponseTitle(err *errorlib.ResponseError) string {
	// enroll a error title bases on response error reason
	return ""
}

func multiErrorCause(err *errorlib.MultiErrors) RootCause {
	// Interpret the error and provide possible cause to fix it
	return RootCause{}
}

func multiErrorResponseTitle(err *errorlib.MultiErrors) string {
	// enroll a error title bases on multi error reason
	return ""
}
