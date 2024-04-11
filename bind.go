package espresso

import (
	"fmt"
	"reflect"
	"strings"
)

// BindSource describes the type of bind.
type BindSource int

const (
	BindPathParam BindSource = iota
	BindFormParam
	BindQueryParam
	BindHeadParam
)

func (b BindSource) String() string {
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
	return fmt.Sprintf("unknown(%d)", int(b))
}

func (b BindSource) Valid() bool {
	return !strings.HasPrefix(b.String(), "unknown")
}

type BindFunc func(any, string) error

type BindParam struct {
	Key  string
	From BindSource
	Type reflect.Type
	Func BindFunc
}

// BindError describes the error when binding a param.
type BindError struct {
	Key  string
	From BindSource
	Type reflect.Type
	Err  error
}

func errorBind(bind BindParam, err error) BindError {
	return BindError{
		Key:  bind.Key,
		From: bind.From,
		Type: bind.Type,
		Err:  err,
	}
}

func (b BindError) Error() string {
	return fmt.Sprintf("bind %s with name %q to type %s error: %v", b.From, b.Key, b.Type, b.Err)
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

func newBindParam(key string, src BindSource, v any) (BindParam, error) {
	vt, fn := getBindFunc(v)
	if fn == nil {
		return BindParam{}, fmt.Errorf("not support to bind %s key %q to %T", src, key, v)
	}

	if !src.Valid() {
		return BindParam{}, fmt.Errorf("not support bind type %d", src)
	}

	return BindParam{
		Key:  key,
		From: src,
		Type: vt,
		Func: fn,
	}, nil
}
