package espresso

import (
	"context"
	"encoding/json"
	"io"
	"mime"

	"github.com/googollee/go-espresso/module"
	"gopkg.in/yaml.v3"
)

var (
	CodecsModule  = module.New[*Codecs]()
	ProvideCodecs = CodecsModule.ProvideWithFunc(func(context.Context) (*Codecs, error) {
		return NewCodecs(JSON{}, YAML{}), nil
	})
)

type Codec interface {
	Mime() string
	Encode(ctx context.Context, w io.Writer, v any) error
	Decode(ctx context.Context, r io.Reader, v any) error
}

type Codecs struct {
	fallback Codec
	codecs   map[string]Codec
}

func NewCodecs(codec ...Codec) *Codecs {
	ret := &Codecs{
		fallback: codec[0],
		codecs:   make(map[string]Codec),
	}

	for _, c := range codec {
		ret.codecs[c.Mime()] = c
	}

	return ret
}

func (c *Codecs) DecodeRequest(ctx Context, v any) error {
	return c.Request(ctx).Decode(ctx, ctx.Request().Body, v)
}

func (c *Codecs) EncodeResponse(ctx Context, v any) error {
	codec := c.Response(ctx)
	return codec.Encode(ctx, ctx.ResponseWriter(), v)
}

func (c *Codecs) Request(ctx Context) Codec {
	ret := c.getCodec(ctx, "Context-Type")
	if ret == nil {
		return c.fallback
	}

	return ret
}

func (c *Codecs) Response(ctx Context) Codec {
	if ret := c.getCodec(ctx, "Accept"); ret != nil {
		return ret
	}

	if ret := c.getCodec(ctx, "Context-Type"); ret != nil {
		return ret
	}

	return c.fallback
}

func (c *Codecs) getCodec(ctx Context, head string) Codec {
	req := ctx.Request()
	reqMime, _, err := mime.ParseMediaType(req.Header.Get(head))
	if err != nil {
		return nil
	}

	ret, ok := c.codecs[reqMime]
	if !ok {
		return nil
	}

	return ret
}

type JSON struct{}

func (JSON) Mime() string {
	return "application/json"
}

func (JSON) Decode(ctx context.Context, r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

func (JSON) Encode(ctx context.Context, w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

type YAML struct{}

func (YAML) Mime() string {
	return "application/yaml"
}

func (YAML) Decode(ctx context.Context, r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}

func (YAML) Encode(ctx context.Context, w io.Writer, v any) error {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()

	return encoder.Encode(v)
}
