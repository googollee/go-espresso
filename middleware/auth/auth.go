package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/googollee/go-expresso/middleware"
)

type AuthWithAkSk[User any] struct {
	bear string
	f    func(ctx context.Context, ak, hash string) (User, error)
	util middleware.Middleware[User]
}

func NewAuthWithAkSk[User any](bear string, f func(ctx context.Context, ak, hash string) (User, error)) *AuthWithAkSk[User] {
	return &AuthWithAkSk[User]{
		bear: bear + " ",
		f:    f,
		util: middleware.NewMiddleware[User](),
	}
}

func (a *AuthWithAkSk[User]) Handle(ctx *gin.Context) {
	auth := ctx.GetHeader("Auth")
	if !strings.HasPrefix(auth, a.bear) {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	auth = strings.Trim(auth[len(a.bear):], "\t ")
	keys := strings.SplitN(auth, ":", 2)
	if len(keys) != 2 {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	ak, hash := keys[0], keys[1]
	user, err := a.f(ctx, ak, hash)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	a.util.Store(ctx, user)
}

func (a *AuthWithAkSk[User]) GetUser(ctx *gin.Context) User {
	return a.util.Get(ctx)
}
