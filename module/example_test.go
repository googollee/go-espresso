package module_test

import (
	"context"
	"fmt"

	"github.com/googollee/go-espresso/module"
)

type DB struct {
	target string
}

func (*DB) CheckHealth(context.Context) error { return nil }

type Cache struct {
	fallback *DB
}

func (*Cache) CheckHealth(context.Context) error { return nil }

func ExampleModule() {
	ctx := context.Background()

	// type DB struct {
	//   target string
	// }
	// func (*DB) CheckHealth(context.Context) error { return nil }
	newDB := func(ctx context.Context) (*DB, error) {
		return &DB{
			target: "localhost.db",
		}, nil
	}
	moduleDB := module.New(newDB)

	// type Cache struct {
	//   fallback *DB
	// }
	// func (*Cache) CheckHealth(context.Context) error { return nil }
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
	// db target: localhost.db
	// cache fallback target: localhost.db
}

func ExampleModule_with_other_value() {
	targetKey := "target"
	ctx := context.WithValue(context.Background(), targetKey, "target.db")

	// type DB struct {
	//   target string
	// }
	// func (*DB) CheckHealth(context.Context) error { return nil }
	newDB := func(ctx context.Context) (*DB, error) {
		target := ctx.Value(targetKey).(string)
		return &DB{
			target: target,
		}, nil
	}
	moduleDB := module.New(newDB)

	// type Cache struct {
	//   fallback *DB
	// }
	// func (*Cache) CheckHealth(context.Context) error { return nil }
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

func ExampleModule_create_with_error() {
	ctx := context.Background()

	// type DB struct {
	//   target string
	// }
	// func (*DB) CheckHealth(context.Context) error { return nil }
	newDB := func(ctx context.Context) (*DB, error) {
		return &DB{
			target: "localhost.db",
		}, nil
	}
	moduleDB := module.New(newDB)

	// type Cache struct {
	//   fallback *DB
	// }
	// func (*Cache) CheckHealth(context.Context) error { return nil }
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

func ExampleModule_create_with_panic() {
	ctx := context.Background()

	// type DB struct {
	//   target string
	// }
	// func (*DB) CheckHealth(context.Context) error { return nil }
	newDB := func(ctx context.Context) (*DB, error) {
		return &DB{
			target: "localhost.db",
		}, nil
	}
	moduleDB := module.New(newDB)

	// type Cache struct {
	//   fallback *DB
	// }
	// func (*Cache) CheckHealth(context.Context) error { return nil }
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

func ExampleModule_no_provider() {
	ctx := context.Background()

	// type DB struct {
	//   target string
	// }
	// func (*DB) CheckHealth(context.Context) error { return nil }
	newDB := func(ctx context.Context) (*DB, error) {
		return &DB{
			target: "localhost.db",
		}, nil
	}
	moduleDB := module.New(newDB)

	// type Cache struct {
	//   fallback *DB
	// }
	// func (*Cache) CheckHealth(context.Context) error { return nil }
	newCache := func(ctx context.Context) (*Cache, error) {
		db := moduleDB.Value(ctx)
		if db == nil {
			return nil, fmt.Errorf("no db as fallback")
		}
		return &Cache{fallback: db}, nil
	}
	moduleCache := module.New(newCache)

	repo := module.NewRepo()
	repo.AddModule(moduleCache)

	_, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	// Output:
	// inject error: module *module_test.Cache creates an instance error: no db as fallback
}
