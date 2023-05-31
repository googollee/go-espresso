package auth

import (
	"net/http"
	"strings"

	"github.com/googollee/go-espresso"
)

type AuthWithAkSk[ContextData any] struct {
	bearer         string
	authParseError error
	authCallback   func(ctx *espresso.Context[ContextData], ak, hash string) error
}

func NewAuthWithAkSk[ContextData any](bearer string, authParseError error, authAndStore func(ctx *espresso.Context[ContextData], ak, hash string) error) *AuthWithAkSk[ContextData] {
	return &AuthWithAkSk[ContextData]{
		bearer:       strings.TrimSpace(bearer) + " ",
		authCallback: authAndStore,
	}
}

func (a *AuthWithAkSk[Data]) Handle(ctx *espresso.Context[Data]) {
	auth := ctx.Request().Header.Get("Auth")
	if !strings.HasPrefix(auth, a.bearer) {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.Abort(a.authParseError)
		return
	}

	auth = strings.Trim(auth[len(a.bearer):], "\t ")
	keys := strings.SplitN(auth, ":", 2)
	if len(keys) != 2 {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.Abort(a.authParseError)
		return
	}

	ak, hash := keys[0], keys[1]
	if err := a.authCallback(ctx, ak, hash); err != nil {
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.Abort(err)
		return
	}
}
