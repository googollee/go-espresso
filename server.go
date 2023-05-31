package espresso

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server[ContextData any] struct {
	server      *http.Server
	router      *httprouter.Router
	initCtxData ContextData
}

func NewServer[ContextData any](init ContextData) *Server[ContextData] {
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

func (s *Server[ContextData]) GET(path string, funcs ...Handler[ContextData]) {
	s.router.GET(path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := newContext(r.Context(), r, w, params, s.initCtxData, funcs)
		ctx.Next()
	})
}

func (s *Server[ContextData]) POST(path string, funcs ...Handler[ContextData]) {
	s.router.POST(path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := newContext(r.Context(), r, w, params, s.initCtxData, funcs)
		ctx.Next()
	})
}
