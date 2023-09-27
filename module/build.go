package module

import (
	"context"

	"github.com/googollee/go-espresso"
)

func Build(ctx context.Context, server *espresso.Server, modules []Module) ([]ModuleImplementer, error) {
	bctx := newBuildContext(ctx, server)

	for _, m := range modules {
		if err := m.build(bctx); err != nil {
			return nil, err
		}
	}

	return bctx.Modules(), nil
}

type buildContext struct {
	context.Context
	server    *espresso.Server
	deps      map[moduleName]struct{}
	instances map[moduleName]ModuleImplementer
	err       error
}

func newBuildContext(ctx context.Context, server *espresso.Server) *buildContext {
	return &buildContext{
		Context:   ctx,
		server:    server,
		deps:      make(map[moduleName]struct{}),
		instances: make(map[moduleName]ModuleImplementer),
	}
}

func (c *buildContext) Child() *buildContext {
	return &buildContext{
		Context:   c.Context,
		server:    c.server,
		deps:      make(map[moduleName]struct{}),
		instances: c.instances,
	}
}

func (c *buildContext) Value(key any) any {
	name, ok := key.(moduleName)
	if !ok {
		return c.Context.Value(key)
	}

	return c.Module(name)
}

func (c *buildContext) Module(name moduleName) ModuleImplementer {
	return c.instances[name]
}

func (c *buildContext) Modules() []ModuleImplementer {
	ret := make([]ModuleImplementer, 0, len(c.instances))
	for _, i := range c.instances {
		ret = append(ret, i)
	}
	return ret
}
