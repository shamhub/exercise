package errors

import (
	"fmt"
	"time"
)

const (
	dateFormat = "2006-01-02T15:04:05.000Z"
)

type CustomError struct {
	StatusCode int         `json:"-" xml:"-"`
	Code       string      `json:"code" xml:"code"`
	Reason     string      `json:"reason" xml:"reason"`
	ResourceID string      `json:"resourceId,omitempty" xml:"resourceId,omitempty"`
	Detail     interface{} `json:"detail,omitempty" xml:"detail,omitempty"`
	Path       string      `json:"path,omitempty" xml:"path,omitempty"`
	RootCauses []RootCause `json:"rootCauses,omitempty" xml:"root_causes,omitempty"`
	Title      string      `json:"title,omitempty" xml:"title,omitempty"`
	DateTime   `json:"datetime" xml:"datetime"`
}

func NewCustomError(code, reason string, statusCode int) *CustomError {
	now := time.Now().UTC()
	formatNow := now.Format(dateFormat)
	return &CustomError{
		Code:       code,
		StatusCode: statusCode,
		Reason:     reason,
		DateTime: DateTime{
			Value:    formatNow,
			TimeZone: "UTC",
		}}
}

type RootCause map[string]interface{}

type DateTime struct {
	Value    string `json:"value" xml:"value"`
	TimeZone string `json:"timezone" xml:"timezone"`
}

func (k *CustomError) Error() string {
	return fmt.Sprint(k.Reason)
}

func (k *CustomError) ErrorResponse(errResp *CustomError) {
	if k != nil {
		*errResp = *k
	}
}
