package espresso

type HTTPCoder interface {
	HTTPCode() int
}

type HTTPIgnore interface {
	Ignore() bool
}

func ErrWithStatus(code int, err error) error {
	return &httpError{
		error:   err,
		Code:    code,
		Message: err.Error(),
	}
}

func ErrWithIgnore(err error) error {
	return httpError{
		error:   err,
		Message: err.Error(),
		ignore:  true,
	}
}

type httpError struct {
	error   `json:"-"`
	Code    int    `json:"-"`
	Message string `json:"message"`

	ignore bool
}

func (e httpError) HTTPCode() int {
	return e.Code
}

func (e httpError) Ignore() bool {
	return e.ignore
}

func (e httpError) Unwrap() error {
	return e.error
}
