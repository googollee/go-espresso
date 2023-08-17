// An Overall Example
// This example shows basic usage of `espresso` framework. It provides
// endpoints to access `Blog` web, as well as APIs to handle `Blog` data. For
// simplicity and focus, `Blog` stores in memory, with a map instance, and all
// code are in one package.
package overall_test

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/googollee/go-espresso"
)

// Definition of `Blog`
type Blog struct {
	ID      int
	Title   string
	Content string
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
		return espresso.ErrWithStatus(http.StatusBadRequest, err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	blog, ok := s.blogs[id]
	if !ok {
		ctx.ResponseWriter().Header().Set("Content-Type", "text/html")
		ctx.ResponseWriter().WriteHeader(http.StatusNotFound)
		fmt.Fprintf(ctx.ResponseWriter(), "<p>not found</p>")
		return espresso.ErrWithStatus(http.StatusNotFound, errors.New("not found"))
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
	if err := ctx.Endpoint(http.MethodPost, "/apis/blog").
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
func OverallExample() {
	server, _ := espresso.New()

	service := &Service{
		blogs: make(map[int]*Blog),
	}

	server.HandleAll(service)
	_ = server.ListenAndServe(":8000")
}
