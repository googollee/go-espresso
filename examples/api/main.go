package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/googollee/go-espresso"
)

type HTTPError struct {
	Code    int    `json:"-"`
	Detail  string `json:"detail"`
	Message string `json:"message"`
}

func (e HTTPError) HTTPCode() int {
	return e.Code
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("(%s)%s", e.Detail, e.Message)
}

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

func (s *Service) Auth(ctx espresso.Context[ContextData]) error {
	auth := ctx.Request().Header.Get("Auth")
	if !strings.HasPrefix(auth, "Bearer user:") {
		return espresso.ErrWithStatus(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	ak := auth[len("Bearer user:"):]
	user, ok := s.users[ak]
	if !ok {
		return espresso.ErrWithStatus(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	ctx.Data().User = user
	return nil
}

type AddArg struct {
	I int
}

type AddReply struct {
	Str string
}

func (s *Service) Add(ctx espresso.Context[ContextData], arg *AddArg) (*AddReply, error) {
	var with int
	if err := ctx.Endpoint(http.MethodPost, "/add/with/:with", s.Auth).
		BindPath("with", &with).
		End(); err != nil {
		return nil, espresso.ErrWithStatus(http.StatusBadRequest, err)
	}

	if arg.I == 0 {
		return nil, &HTTPError{
			Code:    http.StatusBadRequest,
			Detail:  "input_zero",
			Message: "input should not be zero.",
		}
	}

	result := with + arg.I
	ret := AddReply{
		Str: fmt.Sprintf("%d", result),
	}

	return &ret, nil
}

func main() {
	server, err := espresso.NewServer(espresso.WithCodec(espresso.CodecJSON))
	if err != nil {
		log.Fatal("create server error:", err)
	}

	service := &Service{
		users: map[string]*User{
			"access": {
				AccessKey: "access",
				Name:      "user",
			},
		},
	}

	espresso.HandleProcedure(server, ContextData{}, service.Add)

	if err := server.ListenAndServe(":8080"); err != nil {
		log.Fatal("listen and serve error:", err)
	}
}
