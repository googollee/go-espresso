package codec_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/googollee/go-espresso/codec"
)

type Context struct {
	context.Context
	request        *http.Request
	responseWriter http.ResponseWriter
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func ExampleCodecs_jsonContentType() {
	bg := context.Background()
	codecs, err := codec.Default(bg)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://domain/path", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	ctx := &Context{
		Context:        bg,
		request:        req,
		responseWriter: resp,
	}

	reqCodec := codecs.Request(ctx)
	respCodec := codecs.Response(ctx)

	fmt.Println("request:", reqCodec.Mime(), "response:", respCodec.Mime())

	// Output:
	// request: application/json response: application/json
}

func ExampleCodecs_emptyContentType() {
	bg := context.Background()
	codecs, err := codec.Default(bg)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://domain/path", nil)
	if err != nil {
		panic(err)
	}

	resp := httptest.NewRecorder()

	ctx := &Context{
		Context:        bg,
		request:        req,
		responseWriter: resp,
	}

	reqCodec := codecs.Request(ctx)
	respCodec := codecs.Response(ctx)

	fmt.Println("request:", reqCodec.Mime(), "response:", respCodec.Mime())

	// Output:
	// request: application/json response: application/json
}

func ExampleCodecs_differentContextTypeAccept() {
	bg := context.Background()
	codecs, err := codec.Default(bg)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://domain/path", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/yaml")

	resp := httptest.NewRecorder()

	ctx := &Context{
		Context:        bg,
		request:        req,
		responseWriter: resp,
	}

	reqCodec := codecs.Request(ctx)
	respCodec := codecs.Response(ctx)

	fmt.Println("request:", reqCodec.Mime(), "response:", respCodec.Mime())

	// Output:
	// request: application/json response: application/yaml
}
