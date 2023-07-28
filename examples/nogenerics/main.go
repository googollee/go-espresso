//go:build wip

package main

import (
	"context"
	"net/http"

	"github.com/googollee/go-espresso"
)

type User struct {
	ID string
}

func authToUser(user *User) func(ctx espresso.Context) error {
	userID := ctx.Request().Head.Get("Auth")
	user.ID = userID
	return nil
}

type Book struct {
	ID      string
	Name    string
	OwnedBy string
}

type BookService struct{}

func (b *BookService) bindToBook(user *User, book *Book) func(ctx context.Context, id string) error {
	return func(ctx context.Context, id string) error {
		book.ID = id
		book.Name = "existed_book"
		book.OwnedBy = "owner"

		if book.OwnedBy != user.ID {
			return espresso.ErrWithStatus(http.StatusUnauthorized, "unauth")
		}

		return nil
	}
}

func (b *BookService) CreateBook(ctx espresso.Context) error {
	return espresso.Procedure(ctx, b.createBook)
}

func (b *BookService) createBook(ctx espresso.Context, in *Book) (*Book, error) {
	var user User
	if err := ctx.Endpoint(http.MethodPost, "/", authToUser(&user)).
		End(); err != nil {
		return nil, err
	}

	in.ID = "random_id"
	return in, nil
}

func (b *BookService) GetBook(ctx espresso.Context) error {
	return espresso.Provider(ctx, b.getBook)
}

func (b *BookService) getBook(ctx espresso.Context) (*Book, error) {
	var book Book
	var user User
	if err := ctx.Endpoint(http.MethodGet, "/{id}", authToUser(&user)).
		BindPath("id", b.bindToBook(&user, &book)).End(); err != nil {
		return nil, err
	}

	return &book, nil
}

func (b *BookService) UpdateBook(ctx espresso.Context) error {
	return espresso.Procedure(ctx, b.updateBook)
}

func (b *BookService) updateBook(ctx espresso.Context, in *Book) (*Book, error) {
	var org Book
	var user User
	if err := ctx.Endpoint(http.MethodPost, "/{id}", authToUser(&user)).
		BindPath("id", b.bindToBook(&user, &org)).End(); err != nil {
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
	var user User
	if err := ctx.Endpoint(http.MethodDelete, "/{id}", authToUser(&user)).
		BindPath("id", b.bindToBook(&user, &org)).End(); err != nil {
		return nil, err
	}

	return &org, nil
}

func main() {
	svc := BookService{}
	server := espresso.Default()
	server.HandleAll(&svc)
	server.ListenAndServe(":8000")
}
