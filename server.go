package espresso

import (
	"github.com/googollee/go-espresso/module"
)

type Router interface {
	Use(middlewares ...HandleFunc)
	HandleFunc(handleFunc HandleFunc)
	HandleAll(service any)
}

type Server interface {
	Router
	AddModule(provider ...module.Provider)
}

type ServerOption func(s Server) error
