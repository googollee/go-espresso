package runtime

import "github.com/googollee/go-espresso/basetype"

type EndpointBinder struct {
	context  *Context
	endpoint basetype.Endpoint
	err      basetype.BindErrors
}

func (e *EndpointBinder) BindPath(key string, v any, opts ...basetype.BindOption) basetype.EndpointBuilder {
	bindParam := e.endpoint.PathParams[key]
	str := e.context.pathParams.ByName(key)

	if err := bindParam.Func(v, str); err != nil {
		e.err = append(e.err, basetype.ErrBind(bindParam, err))
		return e
	}

	return e
}

func (e *EndpointBinder) BindForm(key string, v any, opts ...basetype.BindOption) basetype.EndpointBuilder {
	bindParam := e.endpoint.FormParams[key]
	req := e.context.request

	if err := req.ParseForm(); err != nil {
		e.err = append(e.err, basetype.ErrBind(bindParam, err))
		return e
	}

	str := req.FormValue(key)
	if err := bindParam.Func(v, str); err != nil {
		e.err = append(e.err, basetype.ErrBind(bindParam, err))
		return e
	}

	return e
}

func (e *EndpointBinder) BindQuery(key string, v any, opts ...basetype.BindOption) basetype.EndpointBuilder {
	bindParam := e.endpoint.QueryParams[key]
	req := e.context.request

	str := req.URL.Query().Get(key)
	if err := bindParam.Func(v, str); err != nil {
		e.err = append(e.err, basetype.ErrBind(bindParam, err))
		return e
	}

	return e
}

func (e *EndpointBinder) BindHead(key string, v any, opts ...basetype.BindOption) basetype.EndpointBuilder {
	bindParam := e.endpoint.HeadParams[key]
	req := e.context.request

	str := req.Header.Get(key)
	if err := bindParam.Func(v, str); err != nil {
		e.err = append(e.err, basetype.ErrBind(bindParam, err))
		return e
	}

	return e
}

func (e *EndpointBinder) End() basetype.BindErrors {
	if len(e.err) != 0 {
		return e.err
	}
	return nil
}
