package espresso

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/julienschmidt/httprouter"
)

type endpoint struct {
	Method       string
	Path         string
	PathParams   []*binding
	FormParams   []*binding
	QueryParams  []*binding
	AcceptMimes  []string
	ResponseMime string
	ResponseType reflect.Type
}

type Declarator interface {
	BindPath(name string, v any) Declarator
	BindForm(name string, v any) Declarator
	BindQuery(name string, v any) Declarator
	Response(mime string) Declarator
	End() BindErrors
}

type endpointBuilder[Data any] struct {
	endpoint   *endpoint
	ctx        *declareContext[Data]
	pathParams map[string]struct{}
}

func (c *declareContext[Data]) Endpoint(method, path string, middleware ...Handler[Data]) Declarator {
	pathParams := make(map[string]struct{})

	router := httprouter.New()
	router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		for _, param := range params {
			pathParams[param.Key] = struct{}{}
		}
	})

	r, err := http.NewRequest(method, path, nil)
	if err != nil {
		err := fmt.Errorf("can't create request with method %s and path %s, it should not happen", method, path)
		panic(err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	c.brew.handlers = middleware

	return &endpointBuilder[Data]{
		endpoint: &endpoint{
			Method: method,
			Path:   path,
		},
		ctx:        c,
		pathParams: pathParams,
	}
}

func (e *endpointBuilder[Data]) BindPath(name string, v any) Declarator {
	if _, ok := e.pathParams[name]; !ok {
		err := fmt.Errorf("can't find variables with name %s in path %s", name, e.endpoint.Path)
		panic(err)
	}
	delete(e.pathParams, name)

	f := getBindFunc(v)
	if f == nil {
		err := fmt.Errorf("can't parse path param %s to type %T", name, v)
		panic(err)
	}

	bind := binding{
		Name:      name,
		BindFunc:  f,
		ValueType: reflect.TypeOf(v).Elem(),
	}
	e.endpoint.PathParams = append(e.endpoint.PathParams, &bind)

	return e
}

func (e *endpointBuilder[Data]) BindForm(name string, v any) Declarator {
	f := getBindFunc(v)
	if f == nil {
		err := fmt.Errorf("can't parse path param %s to type %T", name, v)
		panic(err)
	}

	bind := binding{
		Name:      name,
		BindFunc:  f,
		ValueType: reflect.TypeOf(v).Elem(),
	}
	e.endpoint.FormParams = append(e.endpoint.FormParams, &bind)

	return e
}

func (e *endpointBuilder[Data]) BindQuery(name string, v any) Declarator {
	f := getBindFunc(v)
	if f == nil {
		err := fmt.Errorf("can't parse path param %s to type %T", name, v)
		panic(err)
	}

	bind := binding{
		Name:      name,
		BindFunc:  f,
		ValueType: reflect.TypeOf(v).Elem(),
	}
	e.endpoint.QueryParams = append(e.endpoint.QueryParams, &bind)

	return e
}

func (e *endpointBuilder[Data]) Response(mime string) Declarator {
	e.endpoint.ResponseMime = mime
	return e
}

type endpointDeclareFinished struct{}

func (f endpointDeclareFinished) DeclareDone() bool {
	return true
}

type declareChcecker interface {
	DeclareDone() bool
}

func (e *endpointBuilder[Data]) End() BindErrors {
	if len(e.pathParams) != 0 {
		names := make([]string, 0, len(e.pathParams))
		for name := range e.pathParams {
			names = append(names, name)
		}
		err := fmt.Errorf("didn't bind any variables with path params %v", names)
		panic(err)
	}

	e.ctx.endpoint = e.endpoint

	panic(endpointDeclareFinished{})
}

type handleBinder struct {
	endpoint   *endpoint
	pathIndex  int
	pathParams httprouter.Params
	request    *http.Request
	bindErrors BindErrors
}

func (c *brewContext[Data]) Endpoint(method, path string, handlers ...Handler[Data]) Declarator {
	return &handleBinder{
		endpoint:   c.endpoint,
		pathParams: c.pathParams,
		request:    c.request,
	}
}

func (c *handleBinder) BindPath(name string, v any) Declarator {
	bind := c.endpoint.PathParams[c.pathIndex]
	c.pathIndex++
	if bind.Name != name {
		err := fmt.Errorf("the url param bind is with name %s, should be with name %s", bind.Name, name)
		panic(err)
	}

	if err := bind.BindFunc(c.pathParams.ByName(name), v); err != nil {
		c.bindErrors = append(c.bindErrors, BindError{
			BindType:  BindPathParam,
			ValueType: bind.ValueType,
			Name:      name,
			Err:       err,
		})
		return c
	}

	return c
}

func (c *handleBinder) BindForm(name string, v any) Declarator {
	var bind *binding
	for _, b := range c.endpoint.FormParams {
		if b.Name == name {
			bind = b
			break
		}
	}

	if bind == nil {
		err := fmt.Errorf("can't find a bind of the form param with name %s", name)
		panic(err)
	}

	if c.request.Form == nil {
		const defaultMaxMemory = 32 << 20 // same as http, 32 MB
		c.request.ParseMultipartForm(defaultMaxMemory)
	}
	value := c.request.Form[name]
	if len(value) == 0 {
		return c
	}

	if err := bind.BindFunc(value[0], v); err != nil {
		c.bindErrors = append(c.bindErrors, BindError{
			BindType:  BindFormParam,
			ValueType: bind.ValueType,
			Name:      name,
			Err:       err,
		})
		return c
	}

	return c
}

func (c *handleBinder) BindQuery(name string, v any) Declarator {
	var bind *binding
	for _, b := range c.endpoint.QueryParams {
		if b.Name == name {
			bind = b
			break
		}
	}

	if bind == nil {
		err := fmt.Errorf("can't find a bind of the query param with name %s", name)
		panic(err)
	}

	query := c.request.URL.Query()[name]
	if len(query) == 0 {
		return c
	}

	if err := bind.BindFunc(query[0], v); err != nil {
		c.bindErrors = append(c.bindErrors, BindError{
			BindType:  BindQueryParam,
			ValueType: bind.ValueType,
			Name:      name,
			Err:       err,
		})
		return c
	}

	return c
}

func (c *handleBinder) Response(mime string) Declarator {
	return c
}

func (c *handleBinder) End() BindErrors {
	return c.bindErrors
}
