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

func NewDB(ctx context.Context) (*DB, error) {
	return &DB{
		DB: "db",
	}, nil
}

func (d *DB) CheckHealth(ctx context.Context) error {
	return nil
}

var ModuleDB = module.New[*DB]()

type Cache struct {
	Cache string
}

func NewCache(ctx context.Context) (*Cache, error) {
	db := ModuleDB.Value(ctx)
	return &Cache{
		Cache: "cache with " + db.DB,
	}, nil
}

func (c *Cache) CheckHealth(ctx context.Context) error {
	return nil
}

var ModuleCache = module.New[*Cache]()

type MessageQueue struct {
	Queue string
}

func NewMQ(ctx context.Context) (*MessageQueue, error) {
	return &MessageQueue{
		Queue: "queue",
	}, nil
}

func (r *MessageQueue) CheckHealth(ctx context.Context) error {
	return nil
}

var ModuleMQ = module.New[*MessageQueue]()

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
	server, _ := espresso.Default()

	server.AddModule(ModuleCache.Builder(NewCache), ModuleMQ.Builder(NewMQ), ModuleDB.Builder(NewDB))

	server.HandleFunc(Handler)

	httpSvr := httptest.NewServer(server)
	addr = httpSvr.URL
	cancel = func() {
		httpSvr.Close()
	}

	return
}

func ExampleModule() {
	addr, cancel := LaunchModule()
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
