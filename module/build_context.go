package module

import (
	"context"

	"github.com/googollee/go-espresso"
)

type buildContext struct {
	context.Context
	server  *espresso.Server
	deps    map[moduleName]struct{}
	modules map[moduleName]ModuleImplementer
	err     error
}

func (c *buildContext) Value(key any) any {
	name, ok := key.(moduleName)
	if !ok {
		return c.Context.Value(key)
	}

	return c.modules[name]
}

func (c *buildContext) module(name moduleName) ModuleImplementer {
	return c.modules[name]
}
