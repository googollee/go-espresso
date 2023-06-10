package espresso

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrWithStatus(t *testing.T) {
	const (
		notIgnored = false
		isIgnored  = true
	)
	fakeError := errors.New("error")

	tests := []struct {
		error       error
		wantCode    int
		wantIgnored bool
	}{
		{ErrWithStatus(http.StatusBadRequest, fakeError), http.StatusBadRequest, notIgnored},
		{ErrWithIgnore(fakeError), -1, isIgnored},
	}

	for _, test := range tests {
		t.Run(test.error.Error(), func(t *testing.T) {
			if test.wantIgnored {
				var ig HTTPIgnore
				if errors.As(test.error, &ig) {
					if got, want := ig.Ignore(), true; got != want {
						t.Errorf("(%+v).Ignore() got: %v, want: %v", test.error, got, want)
					}
				} else {
					t.Errorf("error: %+v can't be ignored, want to be ignored", test.error)
				}

			}

			if test.wantCode >= 0 {
				var hc HTTPCoder
				if errors.As(test.error, &hc) {
					if got, want := hc.HTTPCode(), test.wantCode; got != want {
						t.Errorf("(%+v).HTTPCode() got: %v, want: %v", test.error, got, want)
					}
				} else {
					t.Errorf("error: %+v can't get HTTPCode(), want to get HTTPCode()", test.error)
				}
			}

			if !errors.Is(test.error, fakeError) {
				t.Errorf("error %+v should be fakeError, which is not", test.error)
			}
		})
	}
}
