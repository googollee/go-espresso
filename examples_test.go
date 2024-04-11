package espresso_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/googollee/go-espresso"
	"github.com/googollee/go-espresso/codec"
)

type Book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func TestEspresso(t *testing.T) {
	books := make(map[int]Book)
	books[1] = Book{
		ID:    1,
		Title: "The Espresso Book",
	}
	books[2] = Book{
		ID:    2,
		Title: "The Second Book",
	}

	espo := espresso.New()

	espo.HandleFunc(func(ctx espresso.Context) error {
		var id int
		if err := ctx.Endpoint(http.MethodGet, "/book/{id}").
			BindPath("id", &id).
			End(); err != nil {
			return err
		}

		book, ok := books[id]
		if !ok {
			return espresso.Error(http.StatusNotFound, fmt.Errorf("not found"))
		}

		if err := codec.Module.Value(ctx).EncodeResponse(ctx, &book); err != nil {
			return err
		}

		return nil
	})

	espo.HandleFunc(func(ctx espresso.Context) error {
		if err := ctx.Endpoint(http.MethodPost, "/book/").
			End(); err != nil {
			return err
		}

		codecs := codec.Module.Value(ctx)

		var book Book
		if err := codecs.DecodeRequest(ctx, &book); err != nil {
			return espresso.Error(http.StatusBadRequest, err)
		}

		book.ID = len(books)
		books[book.ID] = book

		if err := codecs.EncodeResponse(ctx, &book); err != nil {
			return err
		}

		return nil
	})

	svr := httptest.NewServer(espo)
	defer svr.Close()

	func() {
		var book Book
		resp, err := http.Get(svr.URL + "/book/1")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
			t.Fatal(err)
		}

		if got, want := book.Title, books[1].Title; got != want {
			t.Errorf("got = %q, want: %q", got, want)
		}
	}()

	func() {
		arg := Book{Title: "The New Book"}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(&arg); err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(svr.URL+"/book/1", "application/json", &buf)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		var ret Book
		if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
			panic(err)
		}

		if got, want := ret.ID, 2; got != want {
			t.Errorf("got = %v, want: %v", got, want)
		}
	}()
}
