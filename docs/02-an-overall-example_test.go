// An Overall Example
// This example shows basic usage of `espresso` framework. It provides
// endpoints to access `Blog` web, as well as APIs to handle `Blog` data. For
// simplicity and focus, `Blog` stores in memory, with a map instance, and all
// code are in one package.
package espresso_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/googollee/go-espresso"
)

// Definition of `Blog`
type Blog struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Definition of a service to store `Blog`s
type Service struct {
	mu     sync.RWMutex
	nextID int
	blogs  map[int]*Blog
}

// An endpoint to show a `Blog` as a webpage
func (s *Service) ShowBlogWeb(ctx espresso.Context) error {
	var id int
	if err := ctx.Endpoint(http.MethodGet, "/blogs/:id").
		BindPath("id", &id).
		End(); err != nil {
		ctx.ResponseWriter().Header().Set("Content-Type", "text/html")
		ctx.ResponseWriter().WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(ctx.ResponseWriter(), "<p>bad request</p>")
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	blog, ok := s.blogs[id]
	if !ok {
		ctx.ResponseWriter().Header().Set("Content-Type", "text/html")
		ctx.ResponseWriter().WriteHeader(http.StatusNotFound)
		fmt.Fprintf(ctx.ResponseWriter(), "<p>not found</p>")
		return nil
	}

	ctx.ResponseWriter().Header().Set("Content-Type", "text/html")
	ctx.ResponseWriter().WriteHeader(http.StatusOK)
	fmt.Fprintf(ctx.ResponseWriter(), "<h1>%s</h1><p>%s</p>", blog.Title, blog.Content)

	return nil
}

// An API to create new `Blog`s
func (s *Service) CreateBlog(ctx espresso.Context) error {
	return espresso.Procedure(ctx, s.createBlog)
}

func (s *Service) createBlog(ctx espresso.Context, newBlog *Blog) (*Blog, error) {
	if err := ctx.Endpoint(http.MethodPost, "/api/blogs").
		End(); err != nil {
		return nil, espresso.ErrWithStatus(http.StatusBadRequest, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	newBlog.ID = s.nextID
	s.nextID++
	s.blogs[newBlog.ID] = newBlog

	return newBlog, nil
}

// Register all endpoints and launch the server
func LaunchServer() (addr string, cancel func()) {
	server, _ := espresso.New()

	service := &Service{
		nextID: 1,
		blogs:  make(map[int]*Blog),
	}
	server.HandleAll(service)

	httpSvr := httptest.NewServer(server)
	addr = httpSvr.URL
	cancel = func() {
		httpSvr.Close()
	}

	return
}

func ExampleOverall() {
	addr, cancel := LaunchServer()
	defer cancel()

	// Get a non-exist Blog web
	{
		resp, err := http.Get(addr + "/blogs/1")
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println(resp.StatusCode, resp.Header.Get("Content-Type"), string(body))
	}

	// Get with a bad request
	{
		resp, err := http.Get(addr + "/blogs/non_number")
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println(resp.StatusCode, resp.Header.Get("Content-Type"), string(body))
	}

	// Create a Blog with bad request
	{
		resp, err := http.Post(addr+"/api/blogs", "application/json", strings.NewReader(`
		{
			"invalid_json
		}`))
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println(resp.StatusCode, resp.Header.Get("Content-Type"), strings.TrimSpace(string(body)))
	}

	// Create a Blog
	{
		resp, err := http.Post(addr+"/api/blogs", "application/json", strings.NewReader(`
		{
			"title": "A new web framework",
			"content": "espresso is greate!"
		}`))
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println(resp.StatusCode, resp.Header.Get("Content-Type"), strings.TrimSpace(string(body)))
	}

	// Get the Blog web
	{
		resp, err := http.Get(addr + "/blogs/1")
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println(resp.StatusCode, resp.Header.Get("Content-Type"), string(body))
	}
	// Output:
	// 404 text/html <p>not found</p>
	// 400 text/html <p>bad request</p>
	// 400 application/json {"message":"invalid character '\\n' in string literal"}
	// 200 application/json {"id":1,"title":"A new web framework","content":"espresso is greate!"}
	// 200 text/html <h1>A new web framework</h1><p>espresso is greate!</p>
}
