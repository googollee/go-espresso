package espresso

import (
	"fmt"
	"reflect"
	"strconv"
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
	Type  BindType
	Name  string
	Error error
}

type Binding interface {
	Bind(str string) error
}

type BindErrors []BindError

func (e BindErrors) Error() string {
	return fmt.Sprintf("%v", []BindError(e))
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
