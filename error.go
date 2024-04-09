package espresso

type HTTPError interface {
	HTTPCode() int
}

func Error(code int, err error) error {
	return &httpError{
		Message: err.Error(),

		err:  err,
		code: code,
	}
}

type httpError struct {
	Message string `json:"message"`

	code int
	err  error
}

func (e httpError) Error() string {
	return e.err.Error()
}

func (e httpError) HTTPCode() int {
	return e.code
}

func (e httpError) Unwrap() error {
	return e.err
}
