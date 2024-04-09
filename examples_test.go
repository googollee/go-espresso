package espresso_test

import (
	"fmt"
	"net/http"

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

	espo.Handle(func(ctx espresso.Context) error {
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

		if err := codec.Response(ctx).Encode(ctx, &book); err != nil {
			return err
		}
	})

	espo.Handle(func(ctx espresso.Context) error {
		if err := ctx.Endpoint(http.MethodPost, "/book/").
			End(); err != nil {
			return err
		}

		var book Book
		if err := codec.Request(ctx).Decode(ctx, &book); err != nil {
			return espresso.Error(http.StatusBadRequest, err)
		}

		book.ID = len(books)
		books[book.ID] = book

		if err := codec.Response(ctx).Decode(ctx, &book); err != nil {
			return err
		}
	})
}
