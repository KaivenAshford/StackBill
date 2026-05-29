package service

type ServiceError struct {
	HTTPCode int
	Code     int
	Message  string
}

func NewServiceError(httpCode, code int, message string) *ServiceError {
	return &ServiceError{HTTPCode: httpCode, Code: code, Message: message}
}

func (e *ServiceError) Error() string {
	return e.Message
}
