package errorlib

type HttpResponseError interface {
	error
	ProvideReason() []string
	GetStatusCode() int
}
