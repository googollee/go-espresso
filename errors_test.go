package espresso

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrWithStatus(t *testing.T) {
	fakeError := errors.New("error")

	tests := []struct {
		error    error
		wantCode int
	}{
		{ErrWithStatus(http.StatusBadRequest, fakeError), http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.error.Error(), func(t *testing.T) {
			var hc HTTPCoder
			if !errors.As(test.error, &hc) {
				t.Fatalf("%+v should be a HTTPCoder", test.error)
			}

			if !errors.Is(test.error, fakeError) {
				t.Fatalf("%+v should be a fakeError, which is not", test.error)
			}

			if got, want := hc.HTTPCode(), test.wantCode; got != want {
				t.Errorf("(%+v).HTTPCode() = %v, want: %v", test.error, got, want)
			}
		})
	}
}
