package restapi

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/googollee/go-espresso"
)

type User struct {
	ID       int
	Email    string
	Password string
}

type Blog struct {
	ID       int
	AutherID int
	Title    string
	Content  string
}

type Service struct {
	mu sync.RWMutex

	users map[int]User
	blogs map[int]Blog
}

func NewService(users ...User) *Service {
	ret := &Service{
		users: make(map[int]User),
		blogs: make(map[int]Blog),
	}

	for i, u := range users {
		u.ID = i
		ret.users[u.ID] = u
	}

	return ret
}

func (s *Service) auth(user *User) func(espresso.Context, string) error {
	return func(ctx espresso.Context, authStr string) error {
		sp := strings.SplitN(authStr, ":", 2)
		if len(sp) != 2 {
			return espresso.ErrWithStatus(http.StatusUnauthorized, errors.New("unauth"))
		}

		email, pass := sp[0], sp[1]
		var authUser *User

		for _, u := range s.users {
			if email == u.Email && pass == u.Password {
				authUser = &u
				break
			}
		}

		if authUser == nil {
			return espresso.ErrWithStatus(http.StatusUnauthorized, errors.New("unauth"))
		}

		*user = *authUser
		return nil
	}
}

func (s *Service) bindPathBlog(user User, blog *Blog) func(ctx espresso.Context, id int) error {

	return func(ctx espresso.Context, id int) error {
		b, ok := s.blogs[id]
		if !ok {
			return espresso.ErrWithStatus(http.StatusNotFound, errors.New("not found"))
		}

		if b.AutherID != user.ID {
			return espresso.ErrWithStatus(http.StatusNotFound, errors.New("not your blog"))
		}

		*blog = b
		return nil
	}
}

func (s *Service) CreateBlog(ctx espresso.Context) error {
	return espresso.Produce(ctx, s.createBlog)
}

func (s *Service) createBlog(ctx espresso.Context, input *Blog) (*Blog, error) {
	var user User
	if err := ctx.Endpoint(http.MethodPost, "/blogs").
		BindHead("Authorization", s.auth(&user)).
		End(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	blog := *input
	blog.ID = len(s.blogs)
	blog.AutherID = user.ID

	s.blogs[blog.ID] = blog

	return &blog, nil
}

func (s *Service) GetBlog(ctx espresso.Context) error {
	return espresso.Provide(ctx, s.getBlog)
}

func (s *Service) getBlog(ctx espresso.Context) (*Blog, error) {
	var user User
	var blog Blog
	if err := ctx.Endpoint(http.MethodGet, "/blogs/:id").
		BindHead("Authorization", s.auth(&user)).
		BindPath("id", s.bindPathBlog(user, &blog)).
		End(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return &blog, nil
}

func (s *Service) DeleteBlog(ctx espresso.Context) error {
	return espresso.Provide(ctx, s.deleteBlog)
}

func (s *Service) deleteBlog(ctx espresso.Context) (*Blog, error) {
	var user User
	var blog Blog
	if err := ctx.Endpoint(http.MethodDelete, "/blogs/:id").
		BindHead("Authorization", s.auth(&user)).
		BindPath("id", s.bindPathBlog(user, &blog)).
		End(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	delete(s.blogs, blog.ID)

	return &blog, nil
}
