package module_test

import (
	"context"
	"fmt"

	"github.com/googollee/go-espresso/module"
)

type DB struct {
	target string
}

func NewDB(ctx context.Context) (*DB, error) {
	return &DB{
		target: "localhost.db",
	}, nil
}

var ModuleDB = module.New(NewDB)

type Cache struct {
	fallback *DB
}

func NewCache(ctx context.Context) (*Cache, error) {
	db := ModuleDB.Value(ctx)
	if db == nil {
		return nil, fmt.Errorf("no db as fallback")
	}
	return &Cache{
		fallback: db,
	}, nil
}

var ModuleCache = module.New(NewCache)

func ExampleModule() {
	ctx := context.Background()

	repo := module.NewRepo()
	repo.AddModule(ModuleCache)
	repo.AddModule(ModuleDB)

	ctx, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	db := ModuleDB.Value(ctx)
	cache := ModuleCache.Value(ctx)

	fmt.Println("db target:", db.target)
	fmt.Println("cache fallback target:", cache.fallback.target)

	// Output:
	// db target: localhost.db
	// cache fallback target: localhost.db
}

func ExampleModule_withOtherValue() {
	targetKey := "target"
	ctx := context.WithValue(context.Background(), targetKey, "target.db")

	newDB := func(ctx context.Context) (*DB, error) {
		target := ctx.Value(targetKey).(string)
		return &DB{
			target: target,
		}, nil
	}
	moduleDB := module.New(newDB)

	newCache := func(ctx context.Context) (*Cache, error) {
		db := moduleDB.Value(ctx)
		return &Cache{
			fallback: db,
		}, nil
	}
	moduleCache := module.New(newCache)

	repo := module.NewRepo()
	repo.AddModule(moduleCache)
	repo.AddModule(moduleDB)

	ctx, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	db := moduleDB.Value(ctx)
	cache := moduleCache.Value(ctx)

	_ = db
	_ = cache
	fmt.Println("db target:", db.target)
	fmt.Println("cache fallback target:", cache.fallback.target)

	// Output:
	// db target: target.db
	// cache fallback target: target.db
}

func ExampleModule_createWithError() {
	ctx := context.Background()

	newDB := func(ctx context.Context) (*DB, error) {
		return &DB{
			target: "localhost.db",
		}, nil
	}
	moduleDB := module.New(newDB)

	newCache := func(ctx context.Context) (*Cache, error) {
		_ = moduleDB.Value(ctx)
		return nil, fmt.Errorf("new cache error")
	}
	moduleCache := module.New(newCache)

	repo := module.NewRepo()
	repo.AddModule(moduleCache)
	repo.AddModule(moduleDB)

	_, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	// Output:
	// inject error: module *module_test.Cache creates an instance error: new cache error
}

func ExampleModule_createWithPanic() {
	ctx := context.Background()

	newDB := func(ctx context.Context) (*DB, error) {
		return &DB{
			target: "localhost.db",
		}, nil
	}
	moduleDB := module.New(newDB)

	newCache := func(ctx context.Context) (*Cache, error) {
		_ = moduleDB.Value(ctx)
		panic(fmt.Errorf("new cache error"))
	}
	moduleCache := module.New(newCache)

	repo := module.NewRepo()
	repo.AddModule(moduleCache)
	repo.AddModule(moduleDB)

	defer func() {
		err := recover()
		fmt.Println("panic:", err)
	}()

	_, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	// Output:
	// panic: new cache error
}

func ExampleModule_notExistingProvider() {
	ctx := context.Background()

	repo := module.NewRepo()
	repo.AddModule(ModuleCache)
	// repo.AddModule(moduleDB)

	_, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	// Output:
	// inject error: module *module_test.Cache creates an instance error: no db as fallback
}
