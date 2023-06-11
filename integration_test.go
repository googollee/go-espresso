package espresso

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHandleBaseFunc(t *testing.T) {
	eng, err := NewServer()
	if err != nil {
		t.Fatal("create server error:", err)
	}

	type Data struct {
		User string
	}

	getUser := func(ctx Context[Data]) error {
		user, err := ctx.Request().Cookie("user")
		if err != nil || user.Value == "" {
			return nil
		}

		ctx.Data().User = user.Value
		return nil
	}

	Handle(eng, Data{}, func(ctx Context[Data]) error {
		if err := ctx.Endpoint(http.MethodGet, "/", getUser).
			Response("text/html").
			End(); err != nil {
			return ErrWithStatus(http.StatusBadRequest, err)
		}

		if ctx.Data().User == "" {
			ctx.ResponseWriter().Write([]byte(`<p>Hello from espresso. Please <a href="/login">login</a> first.</p>`))
			return nil
		}

		ret := fmt.Sprintf(`<p>Hello from espresso, %s. Nice to meet you.</p>`, ctx.Data().User)
		ctx.ResponseWriter().Write([]byte(ret))
		return nil
	})

	Handle(eng, Data{}, func(ctx Context[Data]) error {
		if err := ctx.Endpoint(http.MethodGet, "/login").
			Response("text/html").
			End(); err != nil {
			return ErrWithStatus(http.StatusBadRequest, err)
		}

		ctx.ResponseWriter().WriteHeader(http.StatusOK)
		ctx.ResponseWriter().Write([]byte(`<form action="/login">
  <label for="user">user:</label><br>
  <input type="text" id="user" name="user"><br>
  <input type="submit" value="Submit">
</form>`))

		return nil
	})

	Handle(eng, Data{}, func(ctx Context[Data]) error {
		var user string
		if err := ctx.Endpoint(http.MethodPost, "/login").
			Response("text/html").
			BindForm("user", &user).
			End(); err != nil {
			return ErrWithStatus(http.StatusBadRequest, err)
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

	tests := []struct {
		method string
		path   string
		body   io.Reader
		opts   []curlOption

		wantCode int
		wantMime string
		wantBody string
	}{
		{http.MethodGet, "/", nil, nil,
			http.StatusOK, "text/html; charset=utf-8", `<p>Hello from espresso. Please <a href="/login">login</a> first.</p>`},
		{http.MethodGet, "/login", nil, nil,
			http.StatusOK, "text/plain; charset=utf-8", `<form action="/login">
  <label for="user">user:</label><br>
  <input type="text" id="user" name="user"><br>
  <input type="submit" value="Submit">
</form>`},
		{http.MethodPost, "/login", strings.NewReader(url.Values{}.Encode()), []curlOption{withMime("application/x-www-form-urlencoded")},
			http.StatusOK, "text/html; charset=utf-8", `<p>The emtpy user is invalid, Please <a href="/login">login</a>.`},
		{http.MethodPost, "/login", strings.NewReader(url.Values{"user": []string{"my friend"}}.Encode()), []curlOption{withMime("application/x-www-form-urlencoded")},
			http.StatusFound, "", ""},
		{http.MethodGet, "/", nil, []curlOption{withCookie("user", "my friend")},
			http.StatusOK, "text/html; charset=utf-8", `<p>Hello from espresso, my friend. Nice to meet you.</p>`},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("(%s)%s", test.method, test.path), func(t *testing.T) {
			code, mime, resp, err := curl(test.method, baseURL+test.path, test.body, test.opts...)
			if err != nil {
				t.Error("return error:", err)
			}

			if got, want := code, test.wantCode; got != want {
				t.Errorf("code, got: %d, want: %d", got, want)
			}

			if got, want := mime, test.wantMime; got != want {
				t.Errorf("mime, got: %s, want: %s", got, want)
			}

			if diff := cmp.Diff(resp, test.wantBody); diff != "" {
				t.Errorf("body diff: (-got +want)\n%s", diff)
			}
		})
	}
}

type curlOption func(r *http.Request)

func withCookie(name, value string) curlOption {
	return func(r *http.Request) {
		r.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}
}

func withMime(mime string) curlOption {
	return func(r *http.Request) {
		r.Header.Add("Content-Type", mime)
	}
}

func curl(method, url string, bodyReader io.Reader, opts ...curlOption) (code int, mime string, body string, err error) {
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
