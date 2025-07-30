package apperr

type HTTPError struct {
	Code    int
	Message string
}

func (e HTTPError) Error() string {
	return e.Message
}

func New(code int, msg string) HTTPError {
	return HTTPError{Code: code, Message: msg}
}
