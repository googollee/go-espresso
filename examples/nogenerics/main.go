//go:build wip

package main

import (
	"context"
	"net/http"

	"github.com/googollee/go-espresso"
)

type Book struct {
	ID   string
	Name string
}

type BookService struct{}

func (b *BookService) getBook(book *Book) func(ctx context.Context, id string) error {
	return func(ctx context.Context, book *Book) error {
		book.ID = id
		book.Name = "existed_book"
		return nil
	}
}

func (b *BookService) CreateBook(ctx espresso.Context) error {
	return espresso.Procedure(ctx, b.createBook)
}

func (b *BookService) createBook(ctx espresso.Context, in *Book) (*Book, error) {
	if err := ctx.Endpoint(http.MethodPost, "/").
		End(); err != nil {
		return nil, err
	}

	in.ID = "random_id"
	return in, nil
}

func (b *BookService) GetBook(ctx espresso.Context) error {
	var book Book
	if err := ctx.Endpoint(http.MethodGet, "/{id}").
		BindPath("id", b.getBook(&book)).End(); err != nil {
		return nil, err
	}

	return espresso.Provider(ctx, func(context.Context) (*Book, error) {
		return &book, nil
	})
}

func (b *BookService) UpdateBook(ctx espresso.Context) error {
	return espresso.Procedure(ctx, b.updateBook)
}

func (b *BookService) updateBook(ctx espresso.Context, in *Book) (*Book, error) {
	var org Book
	if err := ctx.Endpoint(http.MethodPost, "/{id}").
		BindPath("id", b.getBook(&org)).End(); err != nil {
		return nil, err
	}

	org.Name = in.Name
	return &org, nil
}

func (b *BookService) DeleteBook(ctx espresso.Context) error {
	return espresso.Provider(ctx, b.deleteBook)
}

func (b *BookService) deleteBook(ctx espresso.Context) (*Book, error) {
	var org Book
	if err := ctx.Endpoint(http.MethodDelete, "/{id}").
		BindPath("id", b.getBook(&org)).End(); err != nil {
		return nil, err
	}

	return &org, nil
}

func main() {
	svc := BookService{}
	eng := espresso.Default()
	eng.HandleAll(&svc)
	eng.ListenAndServe(":8000")
}
