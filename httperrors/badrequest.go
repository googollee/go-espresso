package httperrors

import (
	"bytes"
	"fmt"
	"net/http"
)

type FieldError[Code ErrorCodeType] struct {
	Name    string `json:"name"`
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func Field[Code ErrorCodeType](name string, code Code, message string) FieldError[Code] {
	return FieldError[Code]{
		Name:    name,
		Code:    code,
		Message: message,
	}
}

func (e FieldError[Code]) Error() string {
	return fmt.Sprintf("field %s: (%s) %s", e.Name, e.Code, e.Message)
}

type BadRequestError[Code ErrorCodeType] struct {
	Message string                      `json:"message"`
	Details map[string]FieldError[Code] `json:"details"`
}

func BadRequest[Code ErrorCodeType](message string, fields ...FieldError[Code]) *BadRequestError[Code] {
	ret := &BadRequestError[Code]{
		Message: message,
		Details: make(map[string]FieldError[Code]),
	}

	return ret.Join(fields...)
}

func (e *BadRequestError[Code]) Join(fields ...FieldError[Code]) *BadRequestError[Code] {
	for _, field := range fields {
		e.Details[field.Name] = field
	}
	return e
}

func (e BadRequestError[Code]) HTTPCode() int {
	return http.StatusBadRequest
}

func (e BadRequestError[Code]) Error() string {
	var buf bytes.Buffer

	for _, field := range e.Details {
		fmt.Fprintf(&buf, "%s: (%s)%s, ", field.Name, field.Code, field.Message)
	}

	return fmt.Sprintf("(%d)%s: %s", e.HTTPCode(), e.Message, buf.String())
}

func (e BadRequestError[Code]) Unwrap() []error {
	ret := make([]error, 0, len(e.Details))

	for _, field := range e.Details {
		field := field
		ret = append(ret, &field)
	}

	return ret
}
