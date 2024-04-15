package espresso_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/googollee/go-espresso"
)

type Book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func ExampleEspresso() {
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
	// Log to stdout for Output
	espo.AddModule(espresso.LogModule.ProvideWithFunc(func(ctx context.Context) (*slog.Logger, error) {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Remove time from the output for predictable test output.
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		})), nil
	}))
	espo.AddModule(espresso.ProvideCodecs)

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

		if err := espresso.CodecsModule.Value(ctx).EncodeResponse(ctx, &book); err != nil {
			return err
		}

		return nil
	})

	espo.HandleFunc(func(ctx espresso.Context) error {
		if err := ctx.Endpoint(http.MethodPost, "/book/").
			End(); err != nil {
			return err
		}

		codecs := espresso.CodecsModule.Value(ctx)

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
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			panic(resp.Status)
		}

		if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
			panic(err)
		}

		fmt.Println("Book 1 title:", book.Title)
	}()

	func() {
		arg := Book{Title: "The New Book"}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(&arg); err != nil {
			panic(err)
		}

		resp, err := http.Post(svr.URL+"/book/1", "application/json", &buf)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			panic(resp.Status)
		}

		var ret Book
		if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
			panic(err)
		}

		fmt.Println("The New Book id:", ret.ID)
	}()

	// Output:
	// level=INFO msg="receive http" method=GET path=/book/1
	// level=INFO msg="finish http" method=GET path=/book/1
	// Book 1 title: The Espresso Book
	// level=INFO msg="receive http" method=POST path=/book/1
	// level=INFO msg="finish http" method=POST path=/book/1
	// The New Book id: 2
}
