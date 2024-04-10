package module

import (
	"context"
	"testing"
)

func TestProviderReturnNilInterfaceWithoutError(t *testing.T) {
	type Interface interface{}
	newNilInterface := func(context.Context) (Interface, error) {
		return nil, nil
	}
	moduleInterface := New[Interface]()
	provideNil := moduleInterface.ProvideWithFunc(newNilInterface)

	repo := NewRepo()
	repo.AddModule(provideNil)

	ctx, err := repo.InjectTo(context.Background())
	if err != nil {
		t.Fatal("inject error:", err)
	}

	var got Interface = moduleInterface.Value(ctx)
	if got != nil {
		t.Errorf("moduleInterface.Value(ctx) = %v, want: nil", got)
	}
}

func TestProviderReturnNilPointerWithoutError(t *testing.T) {
	type Instance struct{}
	newNilInstance := func(context.Context) (*Instance, error) {
		return nil, nil
	}
	moduleInstance := New[*Instance]()
	provideNil := moduleInstance.ProvideWithFunc(newNilInstance)

	repo := NewRepo()
	repo.AddModule(provideNil)

	ctx, err := repo.InjectTo(context.Background())
	if err != nil {
		t.Fatal("inject error:", err)
	}

	var got *Instance = moduleInstance.Value(ctx)
	if got != nil {
		t.Errorf("moduleInstance.Value(ctx) = %v, want: nil", got)
	}
}

func TestModuleKeyIsNotString(t *testing.T) {
	type Instance struct{}
	newNilInstance := func(context.Context) (*Instance, error) {
		return &Instance{}, nil
	}
	moduleInstance := New[*Instance]()
	provideNil := moduleInstance.ProvideWithFunc(newNilInstance)

	repo := NewRepo()
	repo.AddModule(provideNil)

	ctx, err := repo.InjectTo(context.Background())
	if err != nil {
		t.Fatal("inject error:", err)
	}

	if got := moduleInstance.Value(ctx); got == nil {
		t.Errorf("moduleInstance.Value(ctx) = nil, want: not nil")
	}

	key := moduleInstance.moduleKey
	if got := ctx.Value(key); got == nil {
		t.Errorf("moduleInstance.Value(ctx) = nil, want: not nil")
	}

	keyAsStr := string(key)
	if got := ctx.Value(keyAsStr); got != nil {
		t.Errorf("ctx.Value(string(%q)) = %v, want: nil", keyAsStr, got)
	}
}
