package espresso

import (
	"fmt"
	"strconv"
)

type Binding interface {
	Bind(str string) error
}

type bindFunc func(str string, v any) error

func bind(str string, v any) error {
	b, ok := v.(Binding)
	if !ok {
		panic(fmt.Sprintf("bind(str, v), v is %T, not implements Binding interface"))
	}

	return b.Bind(str)
}

func bindInt(str string, v any) error {
	p, ok := v.(*int)
	if !ok {
		panic(fmt.Sprintf("bindInt(str, v), v is %T, not *int", v))
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}

	*p = int(i)
	return nil
}

func bindStr(str string, v any) error {
	p, ok := v.(*string)
	if !ok {
		panic(fmt.Sprintf("bindStr(str, v), v is %T, not *string", v))
	}

	*p = str
	return nil
}
