package module

import (
	"context"
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

func buildFake[T ModuleImplementer](context.Context, *espresso.Server) (*T, error) {
	var t T
	return &t, nil
}

type Module1 struct{ fakeModule }
type Module2 struct{ fakeModule }
type Module3 struct{ fakeModule }
type Module4 struct{ fakeModule }
type Module5 struct{ fakeModule }

var (
	build1 = buildFake[Module1]
	build2 = buildFake[Module2]
	build3 = buildFake[Module3]
	build4 = buildFake[Module4]
	build5 = buildFake[Module5]
)

var (
	module1 = NewModule(build1)
	module2 = NewModule(build2)
	module3 = NewModule(build3)
	module4 = NewModule(build4)
	module5 = NewModule(build5)
)

func TestModule(t *testing.T) {
	ctx := context.Background()
	server := &espresso.Server{}

	modules, err := Build(ctx, server, []Module{module1, module2, module5})
	if err != nil {
		t.Fatalf("Build(ctx, server, {module1, module2, module5}) returns error: %v", err)
	}

	if got, want := len(modules), 5; got != want {
		t.Fatalf("len(modules) = %v, want: %v", got, want)
	}
}
