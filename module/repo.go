package module

import (
	"context"
	"fmt"
	"maps"
)

// Repo is a repository of modules, and to inject instances creating by modules into a context.
type Repo struct {
	providers map[moduleKey]Provider
}

// NewRepo creates a Repo instance.
func NewRepo() *Repo {
	return &Repo{
		providers: make(map[moduleKey]Provider),
	}
}

// AddModule adds a module to the repo.
// Module always implements Provider, so a module can be added directly.
func (r *Repo) AddModule(provider Provider) {
	r.providers[provider.key()] = provider
}

// InjectTo injects instances created by modules into a context `ctx`.
// It returns a new context with all injections. If any module creates an instance with an error, `InjectTo` returns that error with the module name.
func (r *Repo) InjectTo(ctx context.Context) (ret context.Context, err error) {
	defer func() {
		rErr := recover()
		if rErr == nil {
			return
		}

		createErr, ok := rErr.(createPanic)
		if !ok {
			panic(rErr)
		}

		err = fmt.Errorf("module %s creates an instance error: %w", createErr.key, createErr.err)
	}()

	ret = &moduleContext{
		Context:   ctx,
		providers: maps.Clone(r.providers),
		instances: make(map[moduleKey]any),
	}

	for key := range r.providers {
		_ = ret.Value(key)
	}

	return
}
