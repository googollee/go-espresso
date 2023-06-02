package ctx

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Declarator interface {
	BindPathParam(name string, v any) Declarator
	BindHeader(key string, v any) Declarator
	BindForm(key string, v any) Declarator
	End()
}

type Context[Data any] interface {
	context.Context
	Endpoint(method, path string, middleware ...http.HandlerFunc) Declarator
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Data() Data
}

type AddArg struct {
	I int
}

type AddReply struct {
	Str string
}

type ContextData struct {
	User string
}

type Service struct {
	sessions map[int]*ContextData
}

func (s *Service) LoginPage(ctx Context[ContextData]) {
	ctx.Endpoint(http.MethodGet, "/login").
		Response(http.StatusOK, "text/html").
		End()

	ctx.ResponseWriter().WriteHeader(http.StatusOK)
	ctx.ResponseWriter().Write([]byte(`
<form action="/login">
  <label for="email">email:</label><br>
  <input type="text" id="email" name="email"><br>
  <label for="password">password:</label><br>
  <input type="text" id="password" name="password"><br><br>
  <input type="submit" value="Submit">
</form>`))
}

func (s *Service) Login(ctx Context[ContextData]) {
	var email, password string
	ctx.Endpoint(http.MethodPost, "/login").
		BindForm("email", &email).
		BindForm("password", &password).
		Response(http.StatusUnauthorized, "text/html").
		End()

	if email != "someone@mail.com" || password != "password" {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter().Write([]byte(`
<p>Unaushorized</p>`))
		return
	}

	sessionID := len(s.sessions)
	s.sessions[sessionID] = &ContextData{
		User: "someone@mail.com",
	}
	http.SetCookie(ctx.ResponseWriter(), &http.Cookie{
		Name:     "session",
		Value:    fmt.Sprintf("%d", sessionID),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(ctx.ResponseWriter(), ctx.Request(), "/index.html", http.StatusTemporaryRedirect)
}

func (s *Service) Auth(ctx Context[ContextData]) {
	cookie, err := ctx.Request().Cookie("session")
	if err != nil {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		return
	}

	ses, ok := s.sessions[int(id)]
	if !ok {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx.Data().User = ses.User
}

func (s *Service) Add(ctx Context[struct{}], arg *AddArg) (*AddReply, error) {
	var with int
	var lastModifiedAt time.Time
	ctx.Endpoint(http.MethodPost, "/myservice/add/:with").
		BindPathParam("with", &with).
		BindHeader("Last-Modified-At", &lastModifiedAt).
		End()

	ret := &AddReply{
		Str: fmt.Sprintf("%d", arg.I+with),
	}

	return ret, nil
}
