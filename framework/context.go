package framework

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type Context[T any] struct {
	*gin.Context
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Data           T
}

type RouterGroup[T any] struct {
	eng    *gin.Engine
	prefix string
	init   T
}

func Group[T any](eng *gin.Engine, prefix string, init T) *RouterGroup[T] {
	var v T
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		panic("The type T of Group[T] should be NOT a pointer.")
	}

	return &RouterGroup[T]{
		eng:    eng,
		prefix: prefix,
		init:   init,
	}
}

func (g *RouterGroup[T]) GET(path string, funcs ...func(*Context[T])) {
	handler := func(ctx *gin.Context) {
		ctx_ := &Context[T]{
			Context:        ctx,
			Request:        ctx.Request,
			ResponseWriter: ctx.Writer,
			Data:           g.init,
		}

		for _, fn := range funcs {
			if ctx_.Context.IsAborted() {
				break
			}

			fn(ctx_)
		}
	}

	g.eng.GET(fmt.Sprintf("%s/%s", strings.TrimRight(g.prefix, "/"), strings.TrimLeft(path, "/")), handler)
}
