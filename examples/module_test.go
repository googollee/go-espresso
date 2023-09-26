package espresso_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/googollee/go-espresso"
	"github.com/googollee/go-espresso/module"
)

type DB struct {
	DB string
}

func NewDB(ctx context.Context, s *espresso.Server) (*DB, error) {
	return &DB{
		DB: "db",
	}, nil
}

func (d *DB) CheckHealthy(ctx context.Context) error {
	return nil
}

var ModuleDB = module.NewModule(NewDB)

type Cache struct {
	Cache string
}

func NewCache(ctx context.Context, s *espresso.Server) (*Cache, error) {
	db := ModuleDB.Value(ctx)
	return &Cache{
		Cache: "cache with " + db.DB,
	}, nil
}

func (c *Cache) CheckHealthy(ctx context.Context) error {
	return nil
}

var ModuleCache = module.NewModule(NewCache)

type MessageQueue struct {
	Queue string
}

func NewMQ(ctx context.Context, s *espresso.Server) (*MessageQueue, error) {
	return &MessageQueue{
		Queue: "queue",
	}, nil
}

func (r *MessageQueue) CheckHealthy(ctx context.Context) error {
	return nil
}

var ModuleMQ = module.NewModule(NewMQ)

func Handler(ctx espresso.Context) error {
	if err := ctx.Endpoint(http.MethodGet, "/").End(); err != nil {
		return err
	}

	db := ModuleDB.Value(ctx)
	cache := ModuleCache.Value(ctx)
	mq := ModuleMQ.Value(ctx)

	fmt.Fprintf(ctx.ResponseWriter(), "db: %q, cache: %q, queue: %q", db.DB, cache.Cache, mq.Queue)

	return nil
}

func LaunchModule() (addr string, cancel func()) {
	server, _ := espresso.New()
	if err := server.AddModule(ModuleCache, ModuleMQ, ModuleDB); err != nil {
		panic(err)
	}

	server.HandleFunc(Handler)

	httpSvr := httptest.NewServer(server)
	addr = httpSvr.URL
	cancel = func() {
		httpSvr.Close()
	}

	return
}

func ExampleModule() {
	addr, cancel := LaunchServer()
	defer cancel()

	{
		resp, err := http.Get(addr + "/")
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println(resp.StatusCode, string(body))
	}

	// Output:
	// 200 db: "db", cache: "cache with db", queue: "queue"
}
