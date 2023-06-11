package espresso

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type myValue struct {
	Str string
}

func (v *myValue) Bind(str string) error {
	v.Str = "my:" + str
	return nil
}

type nonBinder struct{}

func TestBind(t *testing.T) {
	type Status int
	const (
		StatusUnknown Status = iota
		StatusOK
		StatusPanic
		StatusError
	)
	tests := []struct {
		v      any
		want   any
		input  string
		status Status
	}{
		{string(""), "string", "string", StatusOK},
		{int(0), 100, "100", StatusOK},
		{myValue{}, myValue{Str: "my:myValue"}, "myValue", StatusOK},

		{int(0), 0, "not a number", StatusError},

		{nonBinder{}, nil, "myValue", StatusPanic},
	}

	for _, test := range tests {
		typ := reflect.TypeOf(test.v)

		t.Run(fmt.Sprintf("%q->type(%s)", test.input, typ.String()), func(t *testing.T) {
			status := StatusUnknown
			defer func() {
				var catch any
				if status == StatusUnknown {
					if catch = recover(); catch != nil {
						status = StatusPanic
					}
				}

				if got, want := status, test.status; got != want {
					t.Errorf("getBindFunc(type %s)(v, %s) got status: %v, want: %v", typ.String(), test.input, got, want)
					if catch != nil {
						t.Errorf("  with panic: %v", catch)
					}
				}
			}()

			value := reflect.New(typ)
			fn := getBindFunc(value.Interface())
			if err := fn(test.input, value.Interface()); err != nil {
				t.Logf("getBindFunc(type %s)(v, %s) got error: %v", typ.String(), test.input, err)
				status = StatusError
				return
			}

			status = StatusOK
			if diff := cmp.Diff(value.Elem().Interface(), test.want); diff != "" {
				t.Errorf("getBindFunc(type %s)(v, %s), v diff: (-got, +want)\n%s", typ.String(), test.input, diff)
			}
		})
	}
}

func TestBindErrors(t *testing.T) {
	fakeError := errors.New("fake error")
	bindErrors := BindErrors{
		{BindType: BindPathParam, ValueType: reflect.TypeOf(int(0)), Name: "path_int", Err: fakeError},
		{BindType: BindFormParam, ValueType: reflect.TypeOf(float64(0)), Name: "form_float", Err: fakeError},
	}

	for _, err := range bindErrors {
		var _ error = err
		if !errors.Is(err, fakeError) {
			t.Errorf("a bind error is not a fakeError, which should be")
		}

		if !strings.Contains(err.Error(), fakeError.Error()) {
			t.Errorf("a bind error is %q, which should contain fakeError %q", err.Error(), fakeError.Error())
		}
	}

	var err error = bindErrors

	if !errors.Is(err, fakeError) {
		t.Errorf("bindErrors is not a fakeError, which should be")
	}

	for _, substr := range []string{
		"path_int", "form_float", fakeError.Error(),
	} {
		if errStr := err.Error(); !strings.Contains(errStr, substr) {
			t.Errorf("a bindErrors is %q, which should contain %q", errStr, substr)
		}
	}
}
