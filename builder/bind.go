package builder

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/googollee/go-espresso/basetype"
)

func newBindParam(key string, src basetype.BindSource, v any) (basetype.BindParam, error) {
	vt, fn := getBindFunc(v)
	if fn == nil {
		return basetype.BindParam{}, fmt.Errorf("not support to bind %s key %q to %T", src, key, v)
	}

	if !src.Valid() {
		return basetype.BindParam{}, fmt.Errorf("not support bind type %d", src)
	}

	return basetype.BindParam{
		Key:  key,
		From: src,
		Type: vt,
		Func: fn,
	}, nil
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func bindInt[T integer](bitSize int) (reflect.Type, basetype.BindFunc) {
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

func bindUint[T uinteger](bitSize int) (reflect.Type, basetype.BindFunc) {
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

func bindFloat[T float](bitSize int) (reflect.Type, basetype.BindFunc) {
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

func bindString[T ~string]() (reflect.Type, basetype.BindFunc) {
	return reflect.TypeOf(T("")), func(v any, param string) error {
		p := v.(*T)
		*p = T(param)
		return nil
	}
}

func getBindFunc(v any) (reflect.Type, basetype.BindFunc) {
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
