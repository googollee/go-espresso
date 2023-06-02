package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/googollee/go-espresso"
)

type User struct {
	AccessKey string
	Name      string
}

type ContextData struct {
	User *User
}

type Service struct {
	users map[string]*User
}

func (s *Service) Auth(ctx *espresso.Context[ContextData]) {
	auth := ctx.Request().Header().Get("Auth")
	if !strings.HasPrefix(auth, "Bearer user:") {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	ak := auth[len("Bearer user:"):]
	user, ok := s.users[ak]
	if !ok {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	ctx.Data.User = user
}

type AddArg struct {
	I int
}

type AddReply struct {
	Str string
}

func (s *Service) Add(ctx *espresso.Context[ContextData], arg *AddArg) (*AddReply, error) {
	var with int
	ctx.Endpoint(http.MethodPost, "/add/with/:with", s.Auth).
		BindPathParam("with", &with).
		End()

	if arg.I == 0 {
		return nil, &HTTPError{
			Code:    http.StatusBadRequest,
			Detail:  "input_zero",
			Message: "input should not be zero.",
		}
	}

	result := with + arg.I
	ret := AddReply{
		Str: fmt.Sprintf("%s", result),
	}

	return ret, nil
}

type HTTPError struct {
	Code    int `json:"-"`
	Detail  string
	Message string
}

func (e HTTPError) HTTPCode() int {
	return e.Code
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("(%s)%s", e.Detail, e.Message)
}

func main() {
	server := espresso.NewServer(ContextData{})

	service := &Service{
		users: map[string]*User{
			"access": &User{
				AccessKey: "access",
				Name:      "user",
			},
		},
	}

	api.Handle(server, service.Add).
		WithDefaultCodec(api.CodecJSON).
		WithErrorType(reflect.TypeOf(HTTPError{}))

	server.ListenAndServe(":8080")
}
