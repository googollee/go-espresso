package espresso

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCodecsJsonContentType(t *testing.T) {
	bg := context.Background()
	codecs := NewCodecs(JSON{}, YAML{})

	req, err := http.NewRequest(http.MethodPost, "http://domain/path", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	ctx := &runtimeContext{
		ctx:      bg,
		request:  req,
		response: resp,
	}

	reqCodec := codecs.Request(ctx)
	if got, want := reqCodec.Mime(), "application/json"; got != want {
		t.Errorf("reqCodec.Mime() = %q, want: %q", got, want)
	}
	respCodec := codecs.Response(ctx)
	if got, want := respCodec.Mime(), "application/json"; got != want {
		t.Errorf("reqCodec.Mime() = %q, want: %q", got, want)
	}
}

func TestCodecsEmptyContentType(t *testing.T) {
	bg := context.Background()
	codecs := NewCodecs(JSON{}, YAML{})

	req, err := http.NewRequest(http.MethodPost, "http://domain/path", nil)
	if err != nil {
		panic(err)
	}

	resp := httptest.NewRecorder()

	ctx := &runtimeContext{
		ctx:      bg,
		request:  req,
		response: resp,
	}

	reqCodec := codecs.Request(ctx)
	if got, want := reqCodec.Mime(), "application/json"; got != want {
		t.Errorf("reqCodec.Mime() = %q, want: %q", got, want)
	}
	respCodec := codecs.Response(ctx)
	if got, want := respCodec.Mime(), "application/json"; got != want {
		t.Errorf("reqCodec.Mime() = %q, want: %q", got, want)
	}
}

func TestCodecsDifferentContextTypeAccept(t *testing.T) {
	bg := context.Background()
	codecs := NewCodecs(JSON{}, YAML{})

	req, err := http.NewRequest(http.MethodPost, "http://domain/path", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/yaml")

	resp := httptest.NewRecorder()

	ctx := &runtimeContext{
		ctx:      bg,
		request:  req,
		response: resp,
	}

	reqCodec := codecs.Request(ctx)
	if got, want := reqCodec.Mime(), "application/json"; got != want {
		t.Errorf("reqCodec.Mime() = %q, want: %q", got, want)
	}
	respCodec := codecs.Response(ctx)
	if got, want := respCodec.Mime(), "application/yaml"; got != want {
		t.Errorf("reqCodec.Mime() = %q, want: %q", got, want)
	}
}
