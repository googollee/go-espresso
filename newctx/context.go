package ctx

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Declarator interface {
	BindPathParam(name string, v any) Declarator
	BindHeader(key string, v any) Declarator
	End()
}

type Context[Data any] interface {
	context.Context
	Endpoint(method, path string, middleware ...http.HandlerFunc) Declarator
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Data() Data
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
	var lastModifiedAt time.Time
	ctx.Endpoint(http.MethodPost, "/myservice/add/:with").
		BindPathParam("with", &with).
		BindHeader("Last-Modified-At", &lastModifiedAt).
		End()

	ret := &AddReply{
		Str: fmt.Sprintf("%d", arg.I+with),
	}

	return ret, nil
}
