package espresso_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/googollee/go-espresso"
)

func ExampleHandle_simpleWeb() {
	eng, err := espresso.NewServer()
	if err != nil {
		log.Fatal("create server error:", err)
	}

	type Data struct {
		User string
	}

	getUser := func(ctx espresso.Context[Data]) error {
		user, err := ctx.Request().Cookie("user")
		if err != nil || user.Value == "" {
			return nil
		}

		ctx.Data().User = user.Value
		return nil
	}

	espresso.Handle(eng, Data{}, func(ctx espresso.Context[Data]) error {
		if err := ctx.Endpoint(http.MethodGet, "/", getUser).
			Response("text/html").
			End(); err != nil {
			return espresso.ErrWithStatus(http.StatusBadRequest, err)
		}

		if ctx.Data().User == "" {
			ctx.ResponseWriter().Write([]byte(`<p>Hello from espresso. Please <a href="/login">login</a> first.</p>`))
			return nil
		}

		ret := fmt.Sprintf(`<p>Hello from espresso, %s. Nice to meet you.</p>`, ctx.Data().User)
		ctx.ResponseWriter().Write([]byte(ret))
		return nil
	})

	espresso.Handle(eng, Data{}, func(ctx espresso.Context[Data]) error {
		if err := ctx.Endpoint(http.MethodGet, "/login").
			Response("text/html").
			End(); err != nil {
			return espresso.ErrWithStatus(http.StatusBadRequest, err)
		}

		ctx.ResponseWriter().WriteHeader(http.StatusOK)
		ctx.ResponseWriter().Write([]byte(`<form action="/login">
  <label for="user">user:</label><br>
  <input type="text" id="user" name="user"><br>
  <input type="submit" value="Submit">
</form>`))

		return nil
	})

	espresso.Handle(eng, Data{}, func(ctx espresso.Context[Data]) error {
		var user string
		if err := ctx.Endpoint(http.MethodPost, "/login").
			Response("text/html").
			BindForm("user", &user).
			End(); err != nil {
			return espresso.ErrWithStatus(http.StatusBadRequest, err)
		}

		if user == "" {
			ctx.ResponseWriter().Write([]byte(`<p>The emtpy user is invalid, Please <a href="/login">login</a>.`))
			return nil
		}

		http.SetCookie(ctx.ResponseWriter(), &http.Cookie{
			Name:  "user",
			Value: user,
			Path:  "/",
		})
		http.Redirect(ctx.ResponseWriter(), ctx.Request(), "/", http.StatusFound)

		return nil
	})

	server := httptest.NewServer(eng)
	defer server.Close()
}
