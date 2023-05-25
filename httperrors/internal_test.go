package httperrors

import (
	"io"
	"testing"
)

func TestInternalAsError(t *testing.T) {
	err := InternalError(io.EOF)
	_, ok := err.(HTTPCoder)
	if got, want := ok, true; got != want {
		t.Errorf("err.(HTTPCoder) = _, %v, want %v", got, want)
	}
}
