package espresso

import (
	"net/http"

	"github.com/googollee/go-espresso/module"
)

type Espresso struct {
	repo   *module.Repo
	mux    *http.ServeMux
	router Router
}

func New() *Espresso {
	ret := &Espresso{
		repo: module.NewRepo(),
		mux:  http.NewServeMux(),
	}
	ret.router = &router{
		mux: ret.mux,
	}

	ret.Use(cacheAllError)

	return ret
}

func (s *Espresso) AddModule(provider ...module.Provider) {
	for _, p := range provider {
		s.repo.Add(p)
	}
}

func (s *Espresso) Use(middlewares ...HandleFunc) {
	s.router.Use(middlewares...)
}

func (s *Espresso) HandleFunc(handleFunc HandleFunc) {
	s.router.HandleFunc(handleFunc)
}

func (s *Espresso) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, err := s.repo.InjectTo(r.Context())
	if err != nil {
		panic(err)
	}

	r = r.WithContext(ctx)
	s.mux.ServeHTTP(w, r)
}
