package module

import (
	"context"
	"fmt"
)

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
		panic(createPanic{key: moduleKey, err: fmt.Errorf("can't find module provideing %q", moduleKey)})
	}

	instance, err := provider.value(c)
	if err != nil {
		panic(createPanic{key: moduleKey, err: err})
	}
	c.instances[moduleKey] = instance
	return instance
}
