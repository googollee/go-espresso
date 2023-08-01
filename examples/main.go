package main

import (
	"github.com/googollee/go-espresso"
	"github.com/googollee/go-espresso/examples/restapi"
)

func main() {
	svr := espresso.New()
	rest := restapi.NewService(
		restapi.User{Email: "person1@domain.com", Password: "123456"},
		restapi.User{Email: "person2@domain.com", Password: "somepass"},
	)

	svr.WithPrefix("/api").HandleAll(rest)

	svr.ListenAndServe(":8000")
}
