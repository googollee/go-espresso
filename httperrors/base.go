package httperrors

import "fmt"

type HTTPCoder interface {
	HTTPCode() int
}

type ErrorCodeType interface {
	~string
}

type DefaultErrorBody[Code ErrorCodeType] struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func (e DefaultErrorBody[Code]) Error() string {
	return fmt.Sprintf("(%s)%s", e.Code, e.Message)
}

type HTTPError struct {
	error
	httpCode int
}

func Error[Code ErrorCodeType](httpCode int, code Code, msg string) error {
	return Err(httpCode, &DefaultErrorBody[Code]{
		Code:    code,
		Message: msg,
	})
}

func Err(httpCode int, err error) error {
	return &HTTPError{
		error:    err,
		httpCode: httpCode,
	}
}

func (e HTTPError) HTTPCode() int {
	return e.httpCode
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("(HTTP %d)%s", e.httpCode, e.error.Error())
}

func (e *HTTPError) Unwrap() error {
	return e.error
}
