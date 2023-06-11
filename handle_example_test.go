package espresso_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

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

type MyHTTPError struct {
	Code    int    `json:"-"`
	Detail  string `json:"detail"`
	Message string `json:"message"`
}

func (e MyHTTPError) HTTPCode() int {
	return e.Code
}

func (e MyHTTPError) Error() string {
	return fmt.Sprintf("(%s)%s", e.Detail, e.Message)
}

func ExampleHandle_restAPI() {
	eng, err := espresso.NewServer(espresso.WithCodec(espresso.CodecJSON))
	if err != nil {
		log.Fatal("create server error:", err)
	}

	type User struct {
		AccessKey string
		Name      string
	}

	type ContextData struct {
		User *User
	}

	users := map[string]*User{
		"access": {
			AccessKey: "access",
			Name:      "name",
		},
	}

	auth := func(ctx espresso.Context[ContextData]) error {
		auth := ctx.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer user:") {
			return &MyHTTPError{
				Code:    http.StatusUnauthorized,
				Detail:  "unauthorized",
				Message: "please add ak",
			}
		}

		ak := auth[len("Bearer user:"):]
		user, ok := users[ak]
		if !ok {
			return &MyHTTPError{
				Code:    http.StatusUnauthorized,
				Detail:  "unauthorized",
				Message: "please add ak",
			}
		}

		ctx.Data().User = user
		return nil
	}

	type AddArg struct {
		I int `json:"i"`
	}

	type AddReply struct {
		Str string `json:"str"`
	}

	espresso.HandleProcedure(eng, ContextData{}, func(ctx espresso.Context[ContextData], arg *AddArg) (*AddReply, error) {
		var with int
		if err := ctx.Endpoint(http.MethodPost, "/add/with/:with", auth).
			BindPath("with", &with).
			End(); err != nil {
			return nil, &MyHTTPError{
				Code:    http.StatusBadRequest,
				Detail:  "bad_request",
				Message: err.Error(),
			}
		}

		result := with + arg.I
		ret := AddReply{
			Str: fmt.Sprintf("%d", result),
		}

		return &ret, nil
	})

	server := httptest.NewServer(eng)
	defer server.Close()

}
