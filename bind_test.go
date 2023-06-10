package espresso

import (
	"fmt"
	"reflect"
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
