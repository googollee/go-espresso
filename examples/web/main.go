package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/googollee/go-espresso"
)

type ContextData struct {
	User string
}

var sessions = make(map[int]*ContextData)

func LoginPage(ctx *espresso.Context[ContextData]) {
	ctx.Endpoint(http.MethodGet, "/login").
		Response("text/html").
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

func Login(ctx *espresso.Context[ContextData]) {
	var email, password string
	ctx.Endpoint(http.MethodPost, "/login").
		BindForm("email", &email).
		BindForm("password", &password).
		Response("text/html").
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

func Auth(ctx *espresso.Context[ContextData]) {
	cookie, err := ctx.Request().Cookie("session")
	if err != nil {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	id, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	ses, ok := sessions[int(id)]
	if !ok {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	ctx.Data = *ses
}

func Index(ctx *espresso.Context[ContextData]) {
	ctx.Endpoint(http.MethodGet, "/index.html", s.Auth).
		Response("text/html").
		End()

	html := fmt.Sprintf("<p>Hello %s from go-espresso</p>", ctx.Data.User)
	ctx.ResponseWriter().Write([]byte(html))
}

func main() {
	server := espresso.NewServer(ContextData{})
	server.Handle(Login, LoginPage, Index)
	server.ListenAndServe(":8080")
}
