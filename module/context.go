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

type moduleContext struct {
	context.Context
	providers map[moduleKey]Provider
	instances map[moduleKey]any
}

func (c *moduleContext) Value(key any) any {
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

	instance, err := provider.value(c)
	if err != nil {
		panic(createPanic{key: moduleKey, err: err})
	}
	c.instances[moduleKey] = instance
	return instance
}
