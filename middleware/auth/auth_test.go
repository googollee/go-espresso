package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/googollee/go-espresso"
)

func TestAuth(t *testing.T) {
	type User struct {
		Name      string
		AccessKey string
	}
	type Data struct {
		User User
	}

	users := map[string]User{
		"user1": {"user1", "ak1"},
		"user2": {"user2", "ak2"},
	}
	auth := NewAuthWithAkSk("Bearer", errors.New("parse auth error"), func(ctx *espresso.Context[Data], ak, hash string) error {
		user, ok := users[ak]
		if !ok {
			return errors.New("not found")
		}

		ctx.Data.User = user

		return nil
	})

	svr := espresso.NewServer(Data{})
	svr.GET("/permission", auth.Handle, func(ctx *espresso.Context[Data]) {
		ctx.ResponseWriter().WriteHeader(http.StatusOK)
	})

	testSvc := httptest.NewServer(svr)
	defer testSvc.Close()

	client := testSvc.Client()

	u := testSvc.URL + "/permission"
	t.Logf("url: %s", u)
	reqWithAuth, _ := http.NewRequest("GET", u, nil)
	reqWithAuth.Header.Add("Auth", "Bearer user1:hash")
	resp, err := client.Do(reqWithAuth)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("response code is %d, want %d", got, want)
	}

	reqNoAuth, _ := http.NewRequest("GET", u, nil)
	resp, err = client.Do(reqNoAuth)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if got, want := resp.StatusCode, http.StatusOK; got == want {
		t.Errorf("response code is %d, want ! %d", got, want)
	}
}
