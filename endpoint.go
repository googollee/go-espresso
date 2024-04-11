package espresso

import "net/http"

type EndpointBuilder interface {
	BindPath(key string, v any) EndpointBuilder
	End() BindErrors
}

type Endpoint struct {
	Method      string
	Path        string
	PathParams  map[string]BindParam
	QueryParams map[string]BindParam
	FormParams  map[string]BindParam
	HeadParams  map[string]BindParam
	ChainFuncs  []HandleFunc
}

func (e *Endpoint) serveHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &runtimeContext{
		ctx:      r.Context(),
		endpoint: e,
		request:  r,
		response: w,
	}
	ctx.Next()
}
