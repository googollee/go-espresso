package espresso

import "context"

type fakeModule struct{}

func (fakeModule) CheckHealthy(context.Context) error { return nil }

var _ Module = &ModuleType[fakeModule]{}
