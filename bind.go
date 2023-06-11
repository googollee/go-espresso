package espresso

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type BindType int

const (
	BindPathParam BindType = iota
	BindFormParam BindType = iota
)

func (b BindType) String() string {
	switch b {
	case BindPathParam:
		return "bind path"
	case BindFormParam:
		return "bind form"
	}
	panic("bind unknown type")
}

type BindError struct {
	BindType  BindType
	ValueType reflect.Type
	Name      string
	Err       error
}

type Binding interface {
	Bind(str string) error
}

func (b BindError) Error() string {
	return fmt.Sprintf("%s with name %s to type %s error: %v", b.BindType, b.Name, b.ValueType, b.Err)
}

func (b BindError) Unwrap() error {
	return b.Err
}

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

type bindFunc func(str string, v any) error

func bindInterface(str string, v any) error {
	b := v.(Binding)
	return b.Bind(str)
}

func bindInt(str string, v any) error {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*v.(*int) = int(i)
	return nil
}

func bindStr(str string, v any) error {
	*v.(*string) = str
	return nil
}

func getBindFunc(v any) bindFunc {
	if _, ok := v.(Binding); ok {
		return bindInterface
	}

	switch v.(type) {
	case *string:
		return bindStr
	case *int:
		return bindInt
	}

	return nil
}

type binding struct {
	Name      string
	BindFunc  bindFunc
	ValueType reflect.Type
}
