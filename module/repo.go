package module

import "context"

type moduleWithInstance struct {
	module    ModuleKey
	dependsOn []contextKey
	instance  Instance
}

type Repo struct {
	mods map[contextKey]moduleWithInstance
}

func NewRepo() *Repo {
	return &Repo{
		mods: make(map[contextKey]moduleWithInstance),
	}
}

func (r *Repo) Add(mods ...ModuleKey) {
	for _, mod := range mods {
		r.mods[mod.contextKey()] = moduleWithInstance{
			module: mod,
		}
	}
}

func (r *Repo) Value(key ModuleKey) Instance {
	m, ok := r.mods[key.contextKey()]
	if !ok {
		return nil
	}

	return m.instance
}

func (r *Repo) Build(ctx context.Context) error {
	var buildNames []contextKey
	for key, mod := range r.mods {
		if mod.instance == nil {
			buildNames = append(buildNames, key)
		}
	}

	bctx := newBuildContext(ctx, r)

	for _, name := range buildNames {
		mod := r.mods[name]
		if err := mod.module.build(bctx); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) addInstance(mod ModuleKey, deps []contextKey, instance Instance) {
	r.mods[mod.contextKey()] = moduleWithInstance{
		module:    mod,
		dependsOn: deps,
		instance:  instance,
	}
}
