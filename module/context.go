package module

import "context"

type errBuildError struct {
	name contextKey
	err  error
}

type buildContext struct {
	context.Context
	mod      ModuleKey
	dependOn map[contextKey]struct{}
	repo     *Repo
}

func newBuildContext(ctx context.Context, repo *Repo) *buildContext {
	return &buildContext{
		Context: ctx,
		repo:    repo,
	}
}

func (ctx *buildContext) Child(mod ModuleKey) *buildContext {
	return &buildContext{
		Context:  ctx.Context,
		mod:      mod,
		dependOn: make(map[contextKey]struct{}),
		repo:     ctx.repo,
	}
}

func (ctx *buildContext) Value(name any) any {
	key, ok := name.(ModuleKey)
	if !ok {
		return ctx.Context.Value(name)
	}

	ctx.dependOn[key.contextKey()] = struct{}{}

	if ret := ctx.repo.Value(key); ret != nil {
		return ret
	}

	if err := key.build(ctx); err != nil {
		panic(errBuildError{name: key.contextKey(), err: err})
	}

	return ctx.repo.Value(key)
}

func (ctx *buildContext) addInstance(instance Instance) {
	deps := make([]contextKey, 0, len(ctx.dependOn))
	for key := range ctx.dependOn {
		deps = append(deps, key)
	}

	ctx.repo.addInstance(ctx.mod, deps, instance)
}
