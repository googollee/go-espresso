package module

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/googollee/assert"
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
	module1 = New[*Module1]()
	module2 = New[*Module2]()
	module3 = New[*Module3]()
	module4 = New[*Module4]()
	module5 = New[*Module5]()
)

func TestRepoBuild(t *testing.T) {
	repo := NewRepo()
	repo.Add(module1.Builder(build1))
	repo.Add(module2.Builder(build2))
	repo.Add(module3.Builder(build3))
	repo.Add(module4.Builder(build4))
	repo.Add(module5.Builder(build5))

	wantModules := []struct {
		module       Key
		wantInstance Instance
	}{
		{module: module1, wantInstance: &Module1{}},
		{module: module2, wantInstance: &Module2{}},
		{module: module3, wantInstance: &Module3{}},
		{module: module4, wantInstance: &Module4{}},
		{module: module5, wantInstance: &Module5{}},
	}

	for _, tc := range wantModules {
		ctx := context.Background()

		got := repo.Value(ctx, tc.module)

		if got, want := fmt.Sprintf("%T", got), fmt.Sprintf("%T", tc.wantInstance); got != want {
			t.Errorf("repo.Value(%v) = %s, want: %s", tc.module, got, want)
		}
	}
}

func TestRepoDependOn(t *testing.T) {
	repo := NewRepo()
	repo.Add(module1.Builder(build1))
	repo.Add(module2.Builder(build2))
	repo.Add(module3.Builder(build3))
	repo.Add(module4.Builder(build4))
	repo.Add(module5.Builder(build5))

	wantDeps := []struct {
		module   Key
		wantDeps assert.Assert[[]key]
	}{
		{module: module1, wantDeps: assert.All(assert.Contain(module2.key), assert.Contain(module3.key))},
		{module: module2, wantDeps: assert.All(assert.Contain(module5.key))},
		{module: module3, wantDeps: assert.Len[key](0)},
		{module: module4, wantDeps: assert.All(assert.Contain(module5.key))},
		{module: module5, wantDeps: assert.Len[key](0)},
	}

	for _, tc := range wantDeps {
		_ = repo.Value(context.Background(), tc.module)

		got := repo.depends[tc.module.name()]
		tc.wantDeps.Checkf(t, got, "repo.depends[%v]", tc.module)
	}
}

func TestRepoBuildError(t *testing.T) {
	buildErr := errors.New("build error")

	panic1 := func(ctx context.Context) (*Module1, error) {
		return nil, fmt.Errorf("module1: %w", buildErr)
	}
	panic2 := func(ctx context.Context) (*Module2, error) {
		return nil, fmt.Errorf("module2: %w", buildErr)
	}

	tests := []struct {
		name          string
		builders      []Builder
		key           Key
		wantErrString string
	}{
		{
			name:          "Direct",
			builders:      []Builder{module1.Builder(panic1)},
			key:           module1,
			wantErrString: "module1",
		},
		{
			name:          "Depend",
			builders:      []Builder{module1.Builder(build1), module2.Builder(panic2)},
			key:           module1,
			wantErrString: "module2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewRepo()
			repo.Add(tc.builders...)

			var err error
			func() {
				defer func() {
					p, ok := recover().(errBuildPanic)
					if !ok {
						return
					}

					err = p.error
				}()

				ctx := context.Background()
				repo.Value(ctx, tc.key)
			}()

			if got, want := err, buildErr; !errors.Is(got, want) {
				t.Fatalf("repo.Build() = %v, want: %v", got, want)
			}
			if got, want := err.Error(), tc.wantErrString; !strings.Contains(got, want) {
				t.Fatalf("repo.Build() = %q, want sub string: %q", got, want)
			}
		})
	}
}

func TestRepoBuildPanic(t *testing.T) {
	buildErr := errors.New("build error")

	panic1 := func(ctx context.Context) (*Module1, error) {
		panic(buildErr)
	}
	panic2 := func(ctx context.Context) (*Module2, error) {
		panic(buildErr)
	}

	tests := []struct {
		name          string
		builders      []Builder
		key           Key
		wantErrString string
	}{
		{
			name:          "Direct",
			builders:      []Builder{module1.Builder(panic1)},
			key:           module1,
			wantErrString: "module1",
		},
		{
			name:          "Depend",
			builders:      []Builder{module1.Builder(build1), module2.Builder(panic2)},
			key:           module1,
			wantErrString: "module1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewRepo()
			repo.Add(tc.builders...)

			func() {
				defer func() {
					r := recover()
					if got, want := r, buildErr; got != want {
						t.Fatalf("repo.Value(ctx, %q) panic with %v, want: %v", tc.key, got, want)
					}
				}()

				ctx := context.Background()
				repo.Value(ctx, tc.key)

				panic("should not reach here")
			}()
		})
	}
}
