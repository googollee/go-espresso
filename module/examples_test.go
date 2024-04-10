package module_test

import (
	"context"
	"fmt"
	"regexp"

	"github.com/googollee/go-espresso/module"
)

type DB interface {
	Target() string
}

type db struct {
	target string
}

func NewDB(ctx context.Context) (DB, error) {
	return &db{
		target: "localhost.db",
	}, nil
}

func (db *db) Target() string {
	return db.target
}

var (
	ModuleDB  = module.New[DB]()
	ProvideDB = ModuleDB.ProvideWithFunc(NewDB)
)

type Cache struct {
	fallback DB
}

func NewCache(ctx context.Context) (*Cache, error) {
	db := ModuleDB.Value(ctx)
	return &Cache{
		fallback: db,
	}, nil
}

var (
	ModuleCache  = module.New[*Cache]()
	ProvideCache = ModuleCache.ProvideWithFunc(NewCache)
)

func ExampleModule() {
	ctx := context.Background()

	repo := module.NewRepo()
	repo.AddModule(ProvideCache)
	repo.AddModule(ProvideDB)

	ctx, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	db := ModuleDB.Value(ctx)
	cache := ModuleCache.Value(ctx)

	fmt.Println("db target:", db.Target())
	fmt.Println("cache fallback target:", cache.fallback.Target())

	// Output:
	// db target: localhost.db
	// cache fallback target: localhost.db
}

func ExampleModule_withOtherValue() {
	targetKey := "target"
	ctx := context.WithValue(context.Background(), targetKey, "target.db")

	repo := module.NewRepo()
	repo.AddModule(ModuleDB.ProvideWithFunc(func(ctx context.Context) (DB, error) {
		target := ctx.Value(targetKey).(string)
		return &db{
			target: target,
		}, nil
	}))
	repo.AddModule(ModuleCache.ProvideWithFunc(func(ctx context.Context) (*Cache, error) {
		db := ModuleDB.Value(ctx)
		return &Cache{
			fallback: db,
		}, nil
	}))

	ctx, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	db := ModuleDB.Value(ctx)
	cache := ModuleCache.Value(ctx)

	fmt.Println("db target:", db.Target())
	fmt.Println("cache fallback target:", cache.fallback.Target())

	// Output:
	// db target: target.db
	// cache fallback target: target.db
}

func ExampleModule_createWithError() {
	ctx := context.Background()

	repo := module.NewRepo()
	repo.AddModule(ModuleDB.ProvideWithFunc(func(ctx context.Context) (DB, error) {
		return &db{
			target: "localhost.db",
		}, nil
	}))
	repo.AddModule(ModuleCache.ProvideWithFunc(func(ctx context.Context) (*Cache, error) {
		_ = ModuleDB.Value(ctx)
		return nil, fmt.Errorf("new cache error")
	}))

	_, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	// Output:
	// inject error: creating with module *module_test.Cache: new cache error
}

func ExampleModule_createWithPanic() {
	ctx := context.Background()

	repo := module.NewRepo()
	repo.AddModule(ModuleDB.ProvideWithFunc(func(ctx context.Context) (DB, error) {
		return &db{
			target: "localhost.db",
		}, nil
	}))
	repo.AddModule(ModuleCache.ProvideWithFunc(func(ctx context.Context) (*Cache, error) {
		_ = ModuleDB.Value(ctx)
		panic(fmt.Errorf("new cache error"))
	}))

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
	repo.AddModule(ProvideCache)
	// repo.AddModule(ProvideDB)

	_, err := repo.InjectTo(ctx)
	if err != nil {
		fmt.Println("inject error:", err)
		return
	}

	// Output:
	// inject error: creating with module module_test.DB: can't find module
}

type FileSystem interface {
	Read()
	Write()
}

type realFileSystem struct{}

func NewRealFileSystem(context.Context) (FileSystem, error) {
	return &realFileSystem{}, nil
}

func (f *realFileSystem) Read()  {}
func (f *realFileSystem) Write() {}

type mockFileSystem struct{}

func NewMockFileSystem(context.Context) (FileSystem, error) {
	return &mockFileSystem{}, nil
}

var (
	ModuleFileSystem      = module.New[FileSystem]()
	ProvideRealFileSystem = ModuleFileSystem.ProvideWithFunc(NewRealFileSystem)
	ProvideMockFileSystem = ModuleFileSystem.ProvideWithFunc(NewMockFileSystem)
)

func (f *mockFileSystem) Read()  {}
func (f *mockFileSystem) Write() {}

func ExampleModule_duplicatingProviders() {
	defer func() {
		p := recover().(string)
		fmt.Println("panic:", regexp.MustCompile(`at .*`).ReplaceAllString(p, "at <removed file and line>"))
	}()

	repo := module.NewRepo()
	repo.AddModule(ProvideRealFileSystem)
	repo.AddModule(ProvideMockFileSystem)

	// Output:
	// panic: already have a provider with type "module_test.FileSystem", added at <removed file and line>
}
