package espresso

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type BindParam struct {
	Key       string
	Type      BindType
	ValueType reflect.Type

	fn bindFunc
}

type BindOption func(b *BindParam) error

func newBindParam(key string, typ BindType, v any) (BindParam, error) {
	vt, fn := getBindFunc(v)
	if fn == nil {
		return BindParam{}, fmt.Errorf("not support to bind %s key %q to %T", typ, key, v)
	}

	if !typ.valid() {
		return BindParam{}, fmt.Errorf("not support bind type %d", typ)
	}

	return BindParam{
		Key:       key,
		Type:      typ,
		ValueType: vt,
		fn:        fn,
	}, nil
}

// BindType describes the type of bind.
type BindType int

const (
	BindPathParam BindType = iota
	BindFormParam
	BindQueryParam
	BindHeadParam
)

func (b BindType) valid() bool {
	return b.String() != "unknown type"
}

func (b BindType) String() string {
	switch b {
	case BindPathParam:
		return "path"
	case BindFormParam:
		return "form"
	case BindQueryParam:
		return "query"
	case BindHeadParam:
		return "head"
	}
	return "unknown type"
}

// BindError describes the error when binding a param.
type BindError struct {
	Key       string
	BindType  BindType
	ValueType reflect.Type
	Err       error
}

func newBindError(b BindParam, err error) BindError {
	return BindError{
		Key:       b.Key,
		BindType:  b.Type,
		ValueType: b.ValueType,
		Err:       err,
	}
}

func (b BindError) Error() string {
	return fmt.Sprintf("bind %s with name %s to type %s error: %v", b.BindType, b.Key, b.ValueType, b.Err)
}

func (b BindError) Unwrap() error {
	return b.Err
}

// BindErrors describes all errors when binding params.
type BindErrors []BindError

func (e BindErrors) Error() string {
	errStr := make([]string, 0, len(e))
	for _, err := range e {
		errStr = append(errStr, err.Error())
	}
	return strings.Join(errStr, ", ")
}

func (e BindErrors) Unwrap() []error {
	if len(e) == 0 {
		return nil
	}

	ret := make([]error, 0, len(e))
	for _, err := range e {
		err := err
		ret = append(ret, err)
	}

	return ret
}

type bindFunc func(any, string) error

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func bindInt[T integer](bitSize int) (reflect.Type, bindFunc) {
	return reflect.TypeOf(T(0)), func(v any, param string) error {
		i, err := strconv.ParseInt(param, 10, bitSize)
		if err != nil {
			return err
		}
		p := v.(*T)
		*p = T(i)
		return nil
	}
}

type uinteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func bindUint[T uinteger](bitSize int) (reflect.Type, bindFunc) {
	return reflect.TypeOf(T(0)), func(v any, param string) error {
		i, err := strconv.ParseUint(param, 10, bitSize)
		if err != nil {
			return err
		}
		p := v.(*T)
		*p = T(i)
		return nil
	}
}

type float interface {
	~float32 | ~float64
}

func bindFloat[T float](bitSize int) (reflect.Type, bindFunc) {
	return reflect.TypeOf(T(0)), func(v any, param string) error {
		i, err := strconv.ParseFloat(param, bitSize)
		if err != nil {
			return err
		}
		p := v.(*T)
		*p = T(i)
		return nil
	}
}

func bindString[T ~string]() (reflect.Type, bindFunc) {
	return reflect.TypeOf(T("")), func(v any, param string) error {
		p := v.(*T)
		*p = T(param)
		return nil
	}
}

func getBindFunc(v any) (reflect.Type, bindFunc) {
	switch v.(type) {
	case *string:
		return bindString[string]()
	case *int:
		return bindInt[int](int(reflect.TypeOf(int(0)).Size()) * 8)
	case *int8:
		return bindInt[int8](8)
	case *int16:
		return bindInt[int16](16)
	case *int32:
		return bindInt[int32](32)
	case *int64:
		return bindInt[int64](64)
	case *uint:
		return bindUint[uint](int(reflect.TypeOf(uint(0)).Size()) * 8)
	case *uint8:
		return bindUint[uint8](8)
	case *uint16:
		return bindUint[uint16](16)
	case *uint32:
		return bindUint[uint32](32)
	case *uint64:
		return bindUint[uint64](64)
	case *float32:
		return bindFloat[float32](32)
	case *float64:
		return bindFloat[float64](64)
	}

	return nil, nil
}
