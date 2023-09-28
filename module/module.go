package module

import (
	"context"
	"errors"
	"fmt"
)

var ErrModuleDependError = errors.New("depend error")
var ErrModuleNotFound = errors.New("not found module")
var errBuildError = errors.New("build error")

type nameKey string

type Module interface {
	Name() nameKey
	DependOn() []nameKey
	CheckHealthy(context.Context) error

	build(*buildContext) error
}

func CheckHealthy(ctx context.Context, names []nameKey) error {
	errs := make(map[nameKey]error)

	checkModuleHealthy(ctx, names, errs)

	if len(errs) != 0 {
		ret := make([]error, 0, len(errs))
		for name, err := range errs {
			ret = append(ret, fmt.Errorf("module %s: %w", name, err))
		}
		return errors.Join(ret...)
	}

	return nil
}

func checkModuleHealthy(ctx context.Context, names []nameKey, errs map[nameKey]error) {
	for _, name := range names {
		if _, ok := errs[name]; ok {
			continue
		}

		v := ctx.Value(name)
		if v == nil {
			errs[name] = ErrModuleNotFound
			continue
		}
		module, ok := v.(Module)
		if !ok {
			errs[name] = ErrModuleNotFound
			continue
		}

		depHealthy := true
		if deps := module.DependOn(); len(deps) != 0 {
			checkModuleHealthy(ctx, deps, errs)

			for _, depname := range deps {
				if errs[depname] != nil {
					depHealthy = false
					break
				}
			}
		}

		if !depHealthy {
			errs[name] = ErrModuleDependError
			continue
		}

		if err := module.CheckHealthy(ctx); err != nil {
			errs[name] = err
		}
	}
}
