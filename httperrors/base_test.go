package httperrors

import (
	"net/http"
	"testing"
)

func TestBaseAsError(t *testing.T) {
	err := Error(http.StatusBadRequest, "bad_request", "reason")
	_, ok := err.(HTTPCoder)
	if got, want := ok, true; got != want {
		t.Errorf("err.(HTTPCoder) = _, %v, want %v", got, want)
	}
}
