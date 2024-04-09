package module

import (
	"context"
	"fmt"
	"runtime"
)

type providerWithLine struct {
	provider Provider
	file     string
	line     int
}

// Repo is a repository of modules, and to inject instances creating by modules into a context.
type Repo struct {
	providers map[moduleKey]providerWithLine
}

// NewRepo creates a Repo instance.
func NewRepo() *Repo {
	return &Repo{
		providers: make(map[moduleKey]providerWithLine),
	}
}

// AddModule adds a module to the repo.
// Module always implements Provider, so a module can be added directly.
func (r *Repo) AddModule(provider Provider) {
	if p, ok := r.providers[provider.key()]; ok {
		msg := fmt.Sprintf("already have a provider with type %q, added at %s:%d", provider.key(), p.file, p.line)
		panic(msg)
	}

	_, file, line, _ := runtime.Caller(1)
	r.providers[provider.key()] = providerWithLine{
		provider: provider,
		file:     file,
		line:     line,
	}
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

	providers := make(map[moduleKey]Provider)
	for k, p := range r.providers {
		providers[k] = p.provider
	}

	ret = &moduleContext{
		Context:   ctx,
		providers: providers,
		instances: make(map[moduleKey]any),
	}

	for key := range r.providers {
		_ = ret.Value(key)
	}

	return
}
