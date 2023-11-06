package basetype

import (
	"github.com/googollee/go-espresso/module"
)

type Router interface {
	Use(middlewares ...HandleFunc)
	HandleFunc(endpoint HandleFunc)
	HandleAll(service any)
}

type Server interface {
	Router
	AddModule(builders ...module.Builder)
}

type ServerOption func(s Server) error
