package module

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/googollee/go-espresso"
)

// - Module1
//   - Module2
//   - Module3
//     - Module5
// - Module2
// - Module4
//   - Module5

type fakeModule struct{}

func (fakeModule) CheckHealthy(context.Context) error { return nil }

type Module1 struct{ fakeModule }
type Module2 struct{ fakeModule }
type Module3 struct{ fakeModule }
type Module4 struct{ fakeModule }
type Module5 struct{ fakeModule }

func build1(ctx context.Context, s *espresso.Server) (*Module1, error) {
	_ = module2.Value(ctx)
	_ = module3.Value(ctx)
	return &Module1{}, nil
}

func build2(ctx context.Context, s *espresso.Server) (*Module2, error) {
	_ = module5.Value(ctx)
	return &Module2{}, nil
}

func build3(ctx context.Context, s *espresso.Server) (*Module3, error) {
	return &Module3{}, nil
}

func build4(ctx context.Context, s *espresso.Server) (*Module4, error) {
	_ = module5.Value(ctx)
	return &Module4{}, nil
}

func build5(ctx context.Context, s *espresso.Server) (*Module5, error) {
	return &Module5{}, nil
}

var (
	module3 = NewModule(build3)
	module5 = NewModule(build5)
	module2 = NewModule(build2)
	module1 = NewModule(build1)
	module4 = NewModule(build4)
)

func TestModule(t *testing.T) {
	ctx := context.Background()
	server := &espresso.Server{}

	modules, err := Build(ctx, server, []Module{module1, module2, module4})
	if err != nil {
		t.Fatalf("Build(ctx, server, {module1, module2, module5}) returns error: %v", err)
	}

	wantModules := map[nameKey]ModuleImplementer{
		module1.Name(): &Module1{},
		module2.Name(): &Module2{},
		module3.Name(): &Module3{},
		module4.Name(): &Module4{},
		module5.Name(): &Module5{},
	}

	if got, want := len(modules), len(wantModules); got != want {
		t.Fatalf("len(modules) = %v, want: %v", got, want)
	}

	for name, wantModule := range wantModules {
		got, ok := modules[name]
		if !ok {
			fmt.Errorf("modules doesn't contain a module with name %q", name)
			continue
		}

		if got, want := fmt.Sprintf("%T", got), fmt.Sprintf("%T", wantModule); got != want {
			fmt.Errorf("modules[%q] = %s, want: %s", name, got, want)
		}
	}

	for mod, wantDeps := range map[Module][]nameKey{
		module1: []nameKey{module2.Name(), module3.Name()},
		module2: nil,
		module3: []nameKey{module5.Name()},
		module4: []nameKey{module5.Name()},
		module5: nil,
	} {
		if got, want := mod.DependOn(), wantDeps; !slices.Equal(got, want) {
			fmt.Errorf("module %q.DependOn() = %v, want: %v", mod.Name(), got, want)
		}
	}
}
