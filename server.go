package espresso

import (
	"net/http"

	"github.com/googollee/go-espresso/module"
)

type Router interface {
	Use(middlewares ...HandleFunc)
	HandleFunc(handleFunc HandleFunc)
	HandleAll(service any)
}

type Server struct {
	repo   *module.Repo
	mux    *http.ServeMux
	router Router
}

func New() *Server {
	return &Server{
		repo:   module.NewRepo(),
		mux:    http.NewServeMux(),
		router: &router{},
	}
}

func (s *Server) AddModule(provider ...module.Provider) {
	for _, p := range provider {
		s.repo.Add(p)
	}
}
