package codec

import (
	"context"
	"mime"

	"github.com/googollee/go-espresso/module"
)

var (
	Module   = module.New[*Codecs]()
	Provider = Module.ProvideWithFunc(Default)
)

type Codecs struct {
	fallback Codec
	codecs   map[string]Codec
}

func Default(context.Context) (*Codecs, error) {
	defaults := []Codec{JSON{}, YAML{}}

	ret := &Codecs{
		fallback: defaults[0],
		codecs:   make(map[string]Codec),
	}

	for _, c := range defaults {
		ret.codecs[c.Mime()] = c
	}

	return ret, nil
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

func (c *Codecs) DecodeRequest(ctx Context, v any) error {
	return c.Request(ctx).Decode(ctx, ctx.Request().Body, v)
}

func (c *Codecs) EncodeResponse(ctx Context, v any) error {
	return c.Response(ctx).Encode(ctx, ctx.ResponseWriter(), v)
}
