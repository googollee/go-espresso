package module

import (
	"context"
	"fmt"
)

type Repo struct {
	instances map[key]Instance
	builders  map[key]Builder
	depends   map[key][]key
}

func NewRepo() *Repo {
	return &Repo{
		instances: make(map[key]Instance),
		builders:  make(map[key]Builder),
		depends:   make(map[key][]key),
	}
}

func (r *Repo) Value(ctx context.Context, k Key) Instance {
	ret, ok := r.instances[k.name()]
	if ok {
		return ret
	}

	builder, ok := r.builders[k.name()]
	if !ok {
		return nil
	}

	bctx := &buildContext{
		Context: ctx,
		repo:    r,
		depends: make(map[key]struct{}),
	}

	ret, err := builder.build(bctx)
	if err != nil {
		panic(errBuildPanic{error: fmt.Errorf("%s build fail: %w", k, err)})
	}

	r.instances[k.name()] = ret

	var deps []key
	if len(bctx.depends) > 0 {
		deps = make([]key, 0, len(bctx.depends))
		for k := range bctx.depends {
			deps = append(deps, k)
		}
	}
	r.depends[k.name()] = deps

	return ret
}

func (r *Repo) Add(builders ...Builder) {
	for _, builder := range builders {
		r.builders[builder.name()] = builder
	}
}
