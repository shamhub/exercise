package errorlib

type HealthCheckError struct {
	statusCode int
	reason     string
}

func NewHealthCheckError(code int, reason string) *HealthCheckError {
	return &HealthCheckError{
		statusCode: code, reason: reason,
	}
}
func (err *HealthCheckError) Error() string {
	return err.reason
}
