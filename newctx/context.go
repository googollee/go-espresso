package ctx

import (
	"context"
	"fmt"
	"net/http"
)

type Context[Data any] interface {
	context.Context
	Endpoint(method, path string, middleware ...http.HandlerFunc)
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Data() Data
	BindPathParam(name string, v any) error
	BindHeader(key string, v any) error
}

type AddArg struct {
	I int
}

type AddReply struct {
	Str string
}

type Service struct{}

func (s *Service) Add(ctx Context[struct{}], arg *AddArg) (*AddReply, error) {
	var with int
	if err := ctx.BindPathParam("with", &with); err != nil {
		return nil, err
	}
	ctx.Endpoint(http.MethodPost, "/myservice/add/:with")

	ret := &AddReply{
		Str: fmt.Sprintf("%d", arg.I+with),
	}

	return ret, nil
}
