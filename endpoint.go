package espresso

import "reflect"

type EndpointBuilder interface {
	BindPath(key string, v any) EndpointBuilder
	End() BindErrors
}

type Endpoint struct {
	Method       string
	Path         string
	PathParams   map[string]BindParam
	QueryParams  map[string]BindParam
	FormParams   map[string]BindParam
	HeadParams   map[string]BindParam
	RequestType  reflect.Type
	ResponseType reflect.Type
	ChainFuncs   []HandleFunc
}

func newEndpoint() *Endpoint {
	return &Endpoint{
		PathParams:  make(map[string]BindParam),
		QueryParams: make(map[string]BindParam),
		FormParams:  make(map[string]BindParam),
		HeadParams:  make(map[string]BindParam),
	}
}
