package module

import (
	"context"
)

func Build(ctx context.Context, modules []Module) (Modules, error) {
	bctx := newBuildContext(ctx)

	for _, m := range modules {
		if err := m.build(bctx); err != nil {
			return nil, err
		}
	}

	return bctx.Modules(), nil
}

type buildContext struct {
	context.Context
	deps      map[nameKey]struct{}
	instances map[nameKey]Instance
	err       error
}

func newBuildContext(ctx context.Context) *buildContext {
	return &buildContext{
		Context:   ctx,
		deps:      make(map[nameKey]struct{}),
		instances: make(map[nameKey]Instance),
	}
}

func (c *buildContext) Child() *buildContext {
	return &buildContext{
		Context:   c.Context,
		deps:      make(map[nameKey]struct{}),
		instances: c.instances,
	}
}

func (c *buildContext) Value(key any) any {
	name, ok := key.(nameKey)
	if !ok {
		return c.Context.Value(key)
	}

	return c.Module(name)
}

func (c *buildContext) Module(name nameKey) Instance {
	return c.instances[name]
}

func (c *buildContext) Modules() map[nameKey]Instance {
	return c.instances
}
