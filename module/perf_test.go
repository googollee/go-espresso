package module

import (
	"context"
	"testing"
)

type DB struct {
	target string
}

func NewDB(ctx context.Context) (*DB, error) {
	return &DB{
		target: "localhost.db",
	}, nil
}

var moduleDB = New(NewDB)

func (*DB) CheckHealth(context.Context) error { return nil }

type Cache struct {
	fallback *DB
}

func NewCache(ctx context.Context) (*Cache, error) {
	db := moduleDB.Value(ctx)
	return &Cache{
		fallback: db,
	}, nil
}

var moduleCache = New(NewCache)

func (*Cache) CheckHealth(context.Context) error { return nil }

func BenchmarkModuleValue(b *testing.B) {
	repo := NewRepo()
	repo.AddModule(moduleDB)
	repo.AddModule(moduleCache)
	ctx, err := repo.InjectTo(context.Background())
	if err != nil {
		b.Fatal("create context error:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = moduleCache.Value(ctx)
	}
}

func BenchmarkContextNomalValue(b *testing.B) {
	ctx := context.Background()

	db, err := NewDB(ctx)
	if err != nil {
		b.Fatal("create db error:", err)
	}
	moduleDB := New(NewDB)
	ctx = context.WithValue(ctx, moduleDB.key(), db)

	cache, err := NewCache(ctx)
	if err != nil {
		b.Fatal("create cache error:", err)
	}
	moduleCache := New(NewCache)
	ctx = context.WithValue(ctx, moduleCache.key(), cache)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Value(moduleCache.key()).(*Cache)
	}
}
