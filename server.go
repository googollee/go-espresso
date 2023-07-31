package espresso

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HandleFunc func(Context) error

type Server struct {
	Group
	endpoints []Endpoint
	router    *httprouter.Router
}

func New() *Server {
	ret := &Server{
		router: httprouter.New(),
	}

	ret.Group = Group{
		prefix: "/",
		server: ret,
	}

	return ret
}

func Default() *Server {
	ret := New()
	ret.Use(
		Logger(
			LogWithMethod(),
			LogWithPath(),
			LogWithMessage(
				func(Context) string { return "received" },
				func(ctx Context) string {
					if ctx.Err() != nil {
						return "error: " + ctx.Err().Error()
					}
					return "succeeded"
				},
			),
		))

	return ret
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
