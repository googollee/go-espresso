package espresso_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/googollee/go-espresso"
)

func ExampleWebServer() {
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

	baseURL := server.URL

	code, mime, resp, err := CURL(http.MethodGet, baseURL+"/", nil)
	fmt.Printf("GET /, response code: %d, mime: %s, err: %v body:\n%s\n\n", code, mime, err, resp)

	code, mime, resp, err = CURL(http.MethodGet, baseURL+"/login", nil)
	fmt.Printf("GET /login, response code: %d, mime: %s, err: %v body:\n%s\n\n", code, mime, err, resp)

	code, mime, resp, err = CURL(http.MethodPost, baseURL+"/login", strings.NewReader(url.Values{}.Encode()), WithMime("application/x-www-form-urlencoded"))
	fmt.Printf("POST /login with nothing, response code: %d, mime: %s, err: %v body:\n%s\n\n", code, mime, err, resp)

	code, mime, resp, err = CURL(http.MethodPost, baseURL+"/login", strings.NewReader(url.Values{"user": []string{"my friend"}}.Encode()), WithMime("application/x-www-form-urlencoded"))
	fmt.Printf("POST /login with user=`my friend`, response code: %d, mime: %s, err: %v body:\n%s\n\n", code, mime, err, resp)

	code, mime, resp, err = CURL(http.MethodGet, baseURL+"/", nil, WithCookie("user", "my friend"))
	fmt.Printf("GET / with cookie, response code: %d, mime: %s, err: %v body:\n%s\n\n", code, mime, err, resp)

	// Output:
	// GET /, response code: 200, mime: text/html, err: <nil> body:
	// <p>Hello from espresso. Please <a href="/login">login</a> first.</p>

	// GET /login, response code: 200, mime: text/html, err: <nil> body:
	// <form action="/login">
	// <label for="user">user:</label><br>
	// <input type="text" id="user" name="user"><br>
	// <input type="submit" value="Submit">
	// </form>

	// POST /login with nothing, response code: 200, mime: text/html, err: <nil> body:
	// <p>The emtpy user is invalid, Please <a href="/login">login</a>.

	// POST /login with user=`my friend`, response code: 302, mime: text/html, err: <nil> body:

	// GET / with cookie, response code: 200, mime: text/html, err: <nil> body:
	// <p>Hello from espresso, my friend. Nice to meet you.</p>
}

type CURLOption func(r *http.Request)

func WithCookie(name, value string) CURLOption {
	return func(r *http.Request) {
		r.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}
}

func WithMime(mime string) CURLOption {
	return func(r *http.Request) {
		r.Header.Add("Content-Type", mime)
	}
}

func CURL(method, url string, bodyReader io.Reader, opts ...CURLOption) (code int, mime string, body string, err error) {
	client := http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return
	}

	for _, opt := range opts {
		opt(req)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	mime = resp.Header.Get("Content-Type")
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	body = string(buf)
	return
}
