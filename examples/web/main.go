package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/googollee/go-espresso"
	"github.com/googollee/go-espresso/httperrors"
)

type ContextData struct {
	User string
}

var (
	sessions  = make(map[int]*ContextData)
	errUnauth = errors.New("please go /login first")
)

func LoginPage(ctx espresso.Context[ContextData]) error {
	if err := ctx.Endpoint(http.MethodGet, "/login").
		Response("text/html").
		End(); err != nil {
		return httperrors.WithStatus(http.StatusBadRequest, err)
	}

	fmt.Println("handle login page")

	ctx.ResponseWriter().WriteHeader(http.StatusOK)
	ctx.ResponseWriter().Write([]byte(`
<form action="/login">
  <label for="email">email:</label><br>
  <input type="text" id="email" name="email"><br>
  <label for="password">password:</label><br>
  <input type="text" id="password" name="password"><br><br>
  <input type="submit" value="Submit">
</form>`))

	return nil
}

func Login(ctx espresso.Context[ContextData]) error {
	var email, password string
	if err := ctx.Endpoint(http.MethodPost, "/login").
		BindForm("email", &email).
		BindForm("password", &password).
		Response("text/html").
		End(); err != nil {
		return httperrors.WithStatus(http.StatusBadRequest, err)
	}

	fmt.Println("handle login with", email, password)

	if email != "someone@mail.com" || password != "password" {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter().Write([]byte(`
<p>Unaushorized</p>`))
		return httperrors.WithStatus(http.StatusUnauthorized, errUnauth)
	}

	sessionID := len(sessions)
	sessions[sessionID] = &ContextData{
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

	return nil
}

func Auth(ctx espresso.Context[ContextData]) error {
	cookie, err := ctx.Request().Cookie("session")
	if err != nil {
		fmt.Println("load cookie error:", err)
		return httperrors.WithStatus(http.StatusUnauthorized, errUnauth)
	}

	id, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		fmt.Println("parse session id", cookie.Value, "error:", err)
		return httperrors.WithStatus(http.StatusUnauthorized, errUnauth)
	}
	fmt.Println("session id", id)

	ses, ok := sessions[int(id)]
	if !ok {
		return httperrors.WithStatus(http.StatusUnauthorized, errUnauth)
	}

	*ctx.Data() = *ses
	return nil
}

func Index(ctx espresso.Context[ContextData]) error {
	if err := ctx.Endpoint(http.MethodGet, "/index.html", Auth).
		Response("text/html").
		End(); err != nil {
		return httperrors.WithStatus(http.StatusBadRequest, err)
	}

	html := fmt.Sprintf("<p>Hello %s from go-espresso</p>", ctx.Data().User)
	ctx.ResponseWriter().Write([]byte(html))

	return nil
}

func main() {
	server := espresso.NewServer(ContextData{})
	server.Handle(Login)
	server.Handle(LoginPage)
	server.Handle(Index)

	fmt.Println("listening with :8080")
	server.ListenAndServe(":8080")
}
