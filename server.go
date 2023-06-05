package espresso

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
)

type Server[ContextData any] struct {
	server      *http.Server
	router      *httprouter.Router
	initCtxData ContextData
}

func NewServer[ContextData any](init ContextData) *Server[ContextData] {
	t := reflect.TypeOf(init)
	if t.Kind() == reflect.Ptr {
		panic("ContextData must NOT be a reference type, nor a pointer.")
	}

	ret := &Server[ContextData]{
		server:      &http.Server{},
		router:      httprouter.New(),
		initCtxData: init,
	}

	ret.server.Handler = ret.router

	return ret
}

func (s *Server[ContextData]) ListenAndServe(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

func (s *Server[ContextData]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server[ContextData]) Handle(fn Handler[ContextData]) {
	declareContext := &declareContext[ContextData]{
		Context: context.Background(),
	}

	func() {
		defer func() {
			r := recover()
			if _, ok := r.(declareChcecker); ok {
				return
			}

			panic(r) // repanic other values.
		}()
		fn(declareContext)
	}()

	endpoint := declareContext.endpoint
	endpoint.Handlers = append(endpoint.Handlers, fn)

	fmt.Println("handle", endpoint.Method, endpoint.Path)
	s.router.Handle(endpoint.Method, endpoint.Path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := handleContext[ContextData]{
			Context:         r.Context(),
			endpoint:        endpoint,
			request:         r,
			responserWriter: w,
			pathParams:      params,
			data:            s.initCtxData,
		}
		ctx.Next()
	})
}
