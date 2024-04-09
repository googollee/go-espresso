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
	moduleInterface := New(newNilInterface)

	repo := NewRepo()
	repo.AddModule(moduleInterface)

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
	moduleInstance := New(newNilInstance)

	repo := NewRepo()
	repo.AddModule(moduleInstance)

	ctx, err := repo.InjectTo(context.Background())
	if err != nil {
		t.Fatal("inject error:", err)
	}

	var got *Instance = moduleInstance.Value(ctx)
	if got != nil {
		t.Errorf("moduleInstance.Value(ctx) = %v, want: nil", got)
	}
}
