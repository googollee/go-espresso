package module

import "context"

type buildContext struct {
	context.Context
	repo    *Repo
	depends map[key]struct{}
}

func (c *buildContext) Value(k any) any {
	if key, ok := k.(key); ok {
		c.depends[key] = struct{}{}
		return c.repo.Value(c, key)
	}

	return c.Context.Value(k)
}
