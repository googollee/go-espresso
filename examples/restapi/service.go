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

func (s *Service) auth(ctx espresso.Context, authStr string) (*User, error) {
	sp := strings.SplitN(authStr, ":", 2)
	if len(sp) != 2 {
		return nil, espresso.ErrWithStatus(http.StatusUnauthorized, errors.New("unauth"))
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
		return nil, espresso.ErrWithStatus(http.StatusUnauthorized, errors.New("unauth"))
	}

	return authUser, nil
}

func (s *Service) getBlogByID(ctx espresso.Context, user *User, id int) (*Blog, error) {
	b, ok := s.blogs[id]
	if !ok {
		return nil, espresso.ErrWithStatus(http.StatusNotFound, errors.New("not found"))
	}

	if b.AutherID != user.ID {
		return nil, espresso.ErrWithStatus(http.StatusNotFound, errors.New("not your blog"))
	}

	return &b, nil
}

func (s *Service) CreateBlog(ctx espresso.Context) error {
	return espresso.Procedure(ctx, s.createBlog)
}

func (s *Service) createBlog(ctx espresso.Context, input *Blog) (*Blog, error) {
	var authStr string
	if err := ctx.Endpoint(http.MethodPost, "/blogs").
		BindHead("Authorization", &authStr).
		End(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, err := s.auth(ctx, authStr)
	if err != nil {
		return nil, err
	}

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
	var authStr string
	var blogID int
	if err := ctx.Endpoint(http.MethodGet, "/blogs/:id").
		BindHead("Authorization", &authStr).
		BindPath("id", &blogID).
		End(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, err := s.auth(ctx, authStr)
	if err != nil {
		return nil, err
	}

	blog, err := s.getBlogByID(ctx, user, blogID)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *Service) DeleteBlog(ctx espresso.Context) error {
	return espresso.Provide(ctx, s.deleteBlog)
}

func (s *Service) deleteBlog(ctx espresso.Context) (*Blog, error) {
	var authStr string
	var blogID int
	if err := ctx.Endpoint(http.MethodDelete, "/blogs/:id").
		BindHead("Authorization", &authStr).
		BindPath("id", &blogID).
		End(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, err := s.auth(ctx, authStr)
	if err != nil {
		return nil, err
	}

	blog, err := s.getBlogByID(ctx, user, blogID)
	if err != nil {
		return nil, err
	}

	delete(s.blogs, blog.ID)

	return blog, nil
}
