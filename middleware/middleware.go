package middleware

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type Middleware[T any] string

func NewMiddleware[T any]() Middleware[T] {
	var t T

	typ := reflect.TypeOf(t)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return Middleware[T](typ.String())
}

func NewMiddlewareWithKey[T any](key string) Middleware[T] {
	return Middleware[T](key)
}

func (m Middleware[T]) Get(ctx *gin.Context) T {
	v, ok := ctx.Keys[string(m)]
	if !ok {
		var t T
		return t
	}

	ret, ok := v.(T)
	if !ok {
		var t T
		return t
	}

	return ret
}

func (m Middleware[T]) Store(ctx *gin.Context, v T) {
	ctx.Keys[string(m)] = v
}
