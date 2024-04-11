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

func newEndpoint() *Endpoint {
	return &Endpoint{
		PathParams:  make(map[string]BindParam),
		QueryParams: make(map[string]BindParam),
		FormParams:  make(map[string]BindParam),
		HeadParams:  make(map[string]BindParam),
	}
}

func (e *Endpoint) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != e.Method {
		http.NotFound(w, r)
		return
	}

	ctx := &runtimeContext{
		ctx:      r.Context(),
		endpoint: e,
		request:  r,
		response: w,
	}
	ctx.Next()
}
