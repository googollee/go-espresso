package espresso

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestBindParam(t *testing.T) {
	tests := []struct {
		key         string
		bindType    BindType
		valueType   reflect.Type
		valueString string

		wantValue any
		wantError string
	}{
		{
			key:         "PathInt",
			bindType:    BindPathParam,
			valueType:   reflect.TypeOf(int(0)),
			valueString: "10",
			wantValue:   int(10),
		},
		{
			key:         "PathUint",
			bindType:    BindPathParam,
			valueType:   reflect.TypeOf(uint(0)),
			valueString: "10",
			wantValue:   uint(10),
		},
		{
			key:         "PathFloat64",
			bindType:    BindPathParam,
			valueType:   reflect.TypeOf(float64(0)),
			valueString: "10.1",
			wantValue:   float64(10.1),
		},
		{
			key:         "PathString",
			bindType:    BindPathParam,
			valueType:   reflect.TypeOf(string("")),
			valueString: "10",
			wantValue:   "10",
		},

		{
			key:         "QueryInt",
			bindType:    BindQueryParam,
			valueType:   reflect.TypeOf(int(0)),
			valueString: "10",
			wantValue:   int(10),
		},
		{
			key:         "QueryUint",
			bindType:    BindQueryParam,
			valueType:   reflect.TypeOf(uint(0)),
			valueString: "10",
			wantValue:   uint(10),
		},
		{
			key:         "QueryFloat64",
			bindType:    BindQueryParam,
			valueType:   reflect.TypeOf(float64(0)),
			valueString: "10.1",
			wantValue:   float64(10.1),
		},
		{
			key:         "QueryString",
			bindType:    BindQueryParam,
			valueType:   reflect.TypeOf(string("")),
			valueString: "10",
			wantValue:   "10",
		},

		{
			key:         "FormInt",
			bindType:    BindFormParam,
			valueType:   reflect.TypeOf(int(0)),
			valueString: "10",
			wantValue:   int(10),
		},
		{
			key:         "FormUint",
			bindType:    BindFormParam,
			valueType:   reflect.TypeOf(uint(0)),
			valueString: "10",
			wantValue:   uint(10),
		},
		{
			key:         "FormFloat64",
			bindType:    BindFormParam,
			valueType:   reflect.TypeOf(float64(0)),
			valueString: "10.1",
			wantValue:   float64(10.1),
		},
		{
			key:         "FormString",
			bindType:    BindFormParam,
			valueType:   reflect.TypeOf(string("")),
			valueString: "10",
			wantValue:   "10",
		},

		{
			key:         "HeadInt",
			bindType:    BindHeadParam,
			valueType:   reflect.TypeOf(int(0)),
			valueString: "10",
			wantValue:   int(10),
		},
		{
			key:         "HeadUint",
			bindType:    BindHeadParam,
			valueType:   reflect.TypeOf(uint(0)),
			valueString: "10",
			wantValue:   uint(10),
		},
		{
			key:         "HeadFloat64",
			bindType:    BindHeadParam,
			valueType:   reflect.TypeOf(float64(0)),
			valueString: "10.1",
			wantValue:   float64(10.1),
		},
		{
			key:         "HeadString",
			bindType:    BindHeadParam,
			valueType:   reflect.TypeOf(string("")),
			valueString: "10",
			wantValue:   "10",
		},

		{
			key:       "InvalidType",
			valueType: reflect.TypeOf(int(0)),
			bindType:  1000,
			wantError: `not support bind type 1000`,
		},
		{
			key:       "InvalidValue",
			valueType: reflect.TypeOf(struct{}{}),
			bindType:  BindPathParam,
			wantError: `not support to bind path key "InvalidValue" to *struct {}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			v := reflect.New(tc.valueType).Interface()

			bindParam, err := newBindParam(tc.key, tc.bindType, v)
			if err != nil || tc.wantError != "" {
				errString := ""
				if err != nil {
					errString = err.Error()
				}
				if got, want := errString, tc.wantError; got != want {
					t.Fatalf("newBindParam(%q, %v, new(%s)) = _, %q, want: %q", tc.key, tc.bindType, tc.valueType.String(), got, want)
				}
				return
			}

			if err := bindParam.fn(v, tc.valueString); err != nil {
				t.Fatalf("newBindParam().fn(new(%s), %q) returns error: %v", tc.valueType.String(), tc.valueString, err)
			}

			if got, want := reflect.ValueOf(v).Elem().Interface(), tc.wantValue; got != want {
				t.Errorf("newBindParam().fn(new(%s), %q) = %v, want: %v", tc.valueType.String(), tc.valueString, got, want)
			}
		})
	}
}

func TestBindError(t *testing.T) {
	underErr := errors.New("my error")
	bErr := newBindError(BindParam{
		Key:       "key",
		Type:      BindPathParam,
		ValueType: reflect.TypeOf(""),
	}, underErr)

	var err error = bErr
	for _, want := range []string{"name \"key\"", "bind path", "type string", "my error"} {
		if got := err.Error(); !strings.Contains(got, want) {
			t.Errorf("bErr.Error() = %q, want substring %q", got, want)
		}
	}

	if !errors.Is(err, underErr) {
		t.Errorf("errors.Is(err, underErr) = false, want: true")
	}
}

func TestBindErrors(t *testing.T) {
	under1 := errors.New("my error1")
	under2 := errors.New("my error2")
	under3 := errors.New("my error3")
	bErrs := BindErrors{
		newBindError(BindParam{
			Key:       "key1",
			Type:      BindPathParam,
			ValueType: reflect.TypeOf(""),
		}, under1),
		newBindError(BindParam{
			Key:       "key2",
			Type:      BindQueryParam,
			ValueType: reflect.TypeOf(1),
		}, under2),
		newBindError(BindParam{
			Key:       "key3",
			Type:      BindFormParam,
			ValueType: reflect.TypeOf(true),
		}, under3),
	}

	var err error = bErrs
	for _, want := range []string{
		"name \"key1\"", "bind path", "type string", "my error1",
		"name \"key2\"", "bind query", "type int", "my error2",
		"name \"key3\"", "bind form", "type bool", "my error3",
	} {
		if got := err.Error(); !strings.Contains(got, want) {
			t.Errorf("bErr.Error() = %q, want substring %q", got, want)
		}
	}

	for _, under := range []error{under1, under2, under3} {
		if !errors.Is(err, under) {
			t.Errorf("errors.Is(err, %v) = false, want: true", under)
		}
	}
}
