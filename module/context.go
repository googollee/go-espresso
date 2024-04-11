package module

import (
	"context"
	"fmt"
)

// ErrNoPrivoder means that it can't find a module within a given context.
// Usually it misses adding that module to a repo.
var ErrNoPrivoder = fmt.Errorf("can't find module")

type createPanic struct {
	key moduleKey
	err error
}

type buildContext struct {
	context.Context
	providers map[moduleKey]providerWithLine
	instances map[moduleKey]any
}

func (c *buildContext) Value(key any) any {
	moduleKey, ok := key.(moduleKey)
	if !ok {
		return c.Context.Value(key)
	}

	if instance, ok := c.instances[moduleKey]; ok {
		return instance
	}

	provider, ok := c.providers[moduleKey]
	if !ok {
		panic(createPanic{key: moduleKey, err: ErrNoPrivoder})
	}

	instance, err := provider.provider.value(c)
	if err != nil {
		panic(createPanic{key: moduleKey, err: err})
	}
	c.instances[moduleKey] = instance
	return instance
}

type moduleContext struct {
	context.Context
	instances map[moduleKey]any
}

func (c *moduleContext) Value(key any) any {
	if moduleKey, ok := key.(moduleKey); ok {
		if instance, ok := c.instances[moduleKey]; ok {
			return instance
		}
	}

	return c.Context.Value(key)
}
