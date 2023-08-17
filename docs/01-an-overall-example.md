# An Overall Example

This example shows basic usage of `espresso` framework. It provides endpoints to access `Blog` web, as well as APIs to handle `Blog` data. For simplicity and focus, `Blog` stores in memory, with a map instance, and all code are in one package.

## Definition of `Blog`

```go
type Blog struct {
    ID int
    Title string
    Content string
}
```

## Definition of a service to store `Blog`s

```go
type Service struct {
    mu sync.RWMutext
    nextID int
    blogs map[int]*Blog
}
```

## An endpoint to show a `Blog` as a webpage

```go
func (s *Service) ShowBlogWeb(ctx espresso.Context) error {
    var id int
    if err := ctx.Endpoint(http.MethodGet, "/blogs/:id").
        BindPath("id", &id).
        End(); err != nil {
        ctx.ResponseWriter().Head().Set("Content-Type", "text/html")
        ctx.ResponseWriter().WriteHead(http.StatusBadRequest)
        fmt.Fprintf("<p>bad request</p>")
        return espresso.ErrWithCode(http.StatusBadRequest, err)
    }

    s.mu.RLock()
    defer s.mu.RUnlock()

    blog, ok := s.blogs[id]
    if !ok {
        ctx.ResponseWriter().Head().Set("Content-Type", "text/html")
        ctx.ResponseWriter().WriteHead(http.StatusNotFound)
        fmt.Fprintf("<p>not found</p>")
        return espresso.ErrWithCode(http.StatusNotFound, errors.New("not found"))
    }

    ctx.ResponseWriter().Head().Set("Content-Type", "text/html")
    ctx.ResponseWriter().WriteHead(http.StatusOK)
    fmt.Fprintf("<h1>%s</h1><p>%s</p>", blog.Title, blog.Content)

    return nil
}
```

## An API to create new `Blog`s

```go
func (s *Service) CreateBlog(ctx espresso.Context) error {
    return espresso.Procedure(s.createBlog)
}

func (s *Service) createBlog(ctx espresso.Context, newBlog *Blog) (*Blog, error) {
    if err := ctx.Endpoint(http.MethodPost, "/apis/blog").
        End(); err != nil {
        return espresso.ErrWithCode(http.StatusBadRequest, err)
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    newBlog.ID= s.nextID
    s.nextID++
    s.blogs[newBlog.ID] = newBlog

    return newBlog, nil
}
```

## Register all endpoints and launch the server

```go
func main() {
    server := espresso.New()

    service := &Service{
        blogs: make(map[int]*Blog),
    }

    server.HandleAll(service)
    server.ListenAndServe(":8000")
}
```
