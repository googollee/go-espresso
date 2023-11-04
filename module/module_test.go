package module

import (
	"context"
	"fmt"
	"slices"
	"testing"
)

// - Module1
//   - Module2
//   - Module3
//     - Module5
// - Module2
// - Module4
//   - Module5

type fakeModule struct{}

func (fakeModule) CheckHealth(context.Context) error { return nil }

type Module1 struct{ fakeModule }
type Module2 struct{ fakeModule }
type Module3 struct{ fakeModule }
type Module4 struct{ fakeModule }
type Module5 struct{ fakeModule }

func build1(ctx context.Context) (*Module1, error) {
	if m := module2.Value(ctx); m == nil {
		panic("module2.Value(ctx) == nil")
	}
	if m := module3.Value(ctx); m == nil {
		panic("module3.Value(ctx) == nil")
	}
	return &Module1{}, nil
}

func build2(ctx context.Context) (*Module2, error) {
	if m := module5.Value(ctx); m == nil {
		panic("module5.Value(ctx) == nil")
	}
	return &Module2{}, nil
}

func build3(ctx context.Context) (*Module3, error) {
	return &Module3{}, nil
}

func build4(ctx context.Context) (*Module4, error) {
	if m := module5.Value(ctx); m == nil {
		panic("module5.Value(ctx) == nil")
	}
	return &Module4{}, nil
}

func build5(ctx context.Context) (*Module5, error) {
	return &Module5{}, nil
}

var (
	module1 = New(build1)
	module2 = New(build2)
	module3 = New(build3)
	module4 = New(build4)
	module5 = New(build5)
)

func TestModule(t *testing.T) {
	ctx := context.Background()

	repo := NewRepo()
	repo.Add(module1)
	repo.Add(module2)
	repo.Add(module4)

	if err := repo.Build(ctx); err != nil {
		t.Fatalf("repo.Build() returns error: %v", err)
	}

	wantModules := []struct {
		module       ModuleKey
		wantInstance Instance
	}{
		{module: module1, wantInstance: &Module1{}},
		{module: module2, wantInstance: &Module2{}},
		{module: module3, wantInstance: &Module3{}},
		{module: module4, wantInstance: &Module4{}},
		{module: module5, wantInstance: &Module5{}},
	}

	for _, tc := range wantModules {
		got := repo.Value(tc.module)

		if got, want := fmt.Sprintf("%T", got), fmt.Sprintf("%T", tc.wantInstance); got != want {
			t.Errorf("repo.Value(%v) = %s, want: %s", tc.module, got, want)
		}
	}
}

func TestModuleDependOn(t *testing.T) {
	ctx := context.Background()

	repo := NewRepo()
	repo.Add(module1)
	repo.Add(module2)
	repo.Add(module4)

	if err := repo.Build(ctx); err != nil {
		t.Fatalf("repo.Build() returns error: %v", err)
	}

	wantDeps := []struct {
		module   ModuleKey
		wantDeps []ModuleKey
	}{
		{module: module1, wantDeps: []ModuleKey{module2, module3}},
		{module: module2, wantDeps: []ModuleKey{module5}},
		{module: module3, wantDeps: nil},
		{module: module4, wantDeps: []ModuleKey{module5}},
		{module: module5, wantDeps: nil},
	}

	for _, tc := range wantDeps {
		var wantDepsKey []contextKey
		for _, mod := range tc.wantDeps {
			wantDepsKey = append(wantDepsKey, mod.contextKey())
		}

		if got, want := repo.mods[tc.module.contextKey()].dependsOn, wantDepsKey; !slices.Equal(got, want) {
			t.Errorf("repo.DependOn(%v) = %v, want: %v", tc.module, got, want)
		}
	}
}
