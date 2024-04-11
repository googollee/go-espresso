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
	repo.Add(provideNil)

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
	repo.Add(provideNil)

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
	newInstance := func(context.Context) (*Instance, error) {
		return &Instance{}, nil
	}
	moduleInstance := New[*Instance]()
	provider := moduleInstance.ProvideWithFunc(newInstance)

	repo := NewRepo()
	repo.Add(provider)

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

func TestSameInstanceBetweenInject(t *testing.T) {
	type Instance struct{}
	newCount := 0
	newInstance := func(context.Context) (*Instance, error) {
		newCount++
		return &Instance{}, nil
	}
	moduleInstance := New[*Instance]()
	provider := moduleInstance.ProvideWithFunc(newInstance)

	repo := NewRepo()
	repo.Add(provider)

	ctx1, err := repo.InjectTo(context.Background())
	if err != nil {
		t.Fatal("inject error:", err)
	}
	instance1 := moduleInstance.Value(ctx1)

	ctx2, err := repo.InjectTo(context.Background())
	if err != nil {
		t.Fatal("inject error:", err)
	}
	instance2 := moduleInstance.Value(ctx2)

	if instance1 != instance2 {
		t.Fatal("different instances.")
	}

	if newCount != 1 {
		t.Fatalf("newInstance was run %d times, which should be only 1.", newCount)
	}
}
