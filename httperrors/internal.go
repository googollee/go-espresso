package httperrors

import "net/http"

func InternalError(err error) error {
	return &HTTPError{
		error:    err,
		httpCode: http.StatusInternalServerError,
	}
}
