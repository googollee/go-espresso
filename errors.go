package espresso

type HTTPCoder interface {
	HTTPCode() int
}

type HTTPIgnore interface {
	Ignore() bool
}

func WithStatus(code int, err error) error {
	return &httpError{
		error: err,
		Code:  code,
	}
}

func WithIgnore(err error) error {
	return &httpIgnore{
		error: err,
	}
}

type httpError struct {
	error
	Code int
}

func (e httpError) HTTPCode() int {
	return e.Code
}

func (e httpError) Unwrap() error {
	return e.error
}

type httpIgnore struct {
	error
}

func (e httpIgnore) Ignore() bool {
	return true
}

func (e httpIgnore) Unwrap() error {
	return e.error
}
