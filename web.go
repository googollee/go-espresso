package espresso

import (
	"fmt"
	"html/template"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
)

type WebOpt func(*Web) error

type Web struct {
	templates *template.Template
}

func NewWeb(tempaltes *template.Template, opts ...WebOpt) (*Web, error) {
	ret := Web{
		templates: tempaltes,
	}

	for _, opt := range opts {
		if err := opt(&ret); err != nil {
			return nil, fmt.Errorf("configure web error: %w", err)
		}
	}

	return &ret, nil
}

func (w *Web) LoadForm(ctx *gin.Context, v any) error {
	if err := ctx.Request.ParseForm(); err != nil {
		return fmt.Errorf("invalid form: %w", err)
	}

	if err := form.NewDecoder().Decode(&v, ctx.Request.Form); err != nil {
		return fmt.Errorf("parse form error: %w", err)
	}

	return nil
}

func (w *Web) Render(ctx *gin.Context, code int, tmpl, mime string, v any) {
	wr := ctx.Writer
	wr.Header().Add("Content-Type", mime)
	wr.WriteHeader(code)

	w.templates.ExecuteTemplate(wr, tmpl, v)
}

func (w *Web) Response(ctx *gin.Context, code int, mime string, r io.Reader) {
	wr := ctx.Writer
	wr.Header().Add("Content-Type", mime)
	wr.WriteHeader(code)

	io.Copy(wr, r)
}
