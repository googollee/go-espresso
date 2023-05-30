package framework

import (
	"net/http"
	"strings"
)

type AuthWithAkSk[CTX any] struct {
	bear string
	fn   func(ctx *Context[CTX], ak, sk string) error
}

func NewAuthWithAkSk[CTX any](bear string, fn func(ctx *Context[CTX], ak, sk string) error) *AuthWithAkSk[CTX] {
	return &AuthWithAkSk[CTX]{
		bear: strings.TrimSpace(bear) + " ",
		fn:   fn,
	}
}

func (a *AuthWithAkSk[CTX]) Handle(ctx *Context[CTX]) {
	auth := ctx.Request.Header.Get("Auth")
	if !strings.HasPrefix(auth, a.bear) {
		ctx.Status(http.StatusNetworkAuthenticationRequired)
		ctx.Abort()
		return
	}

	auth = auth[len(a.bear):]
	keys := strings.SplitN(auth, ":", 2)
	if len(keys) != 2 {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	ak, sk := keys[0], keys[1]
	if err := a.fn(ctx, ak, sk); err != nil {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}
}
