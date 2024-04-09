package codec

import (
	"context"
	"mime"

	"github.com/googollee/go-espresso/module"
)

var Module = module.New(Default)

type Codecs struct {
	fallback Codec
	codecs   map[string]Codec
}

func Default(context.Context) (*Codecs, error) {
	json := JSON{}

	return &Codecs{
		fallback: json,
		codecs: map[string]Codec{
			json.Mime(): json,
		},
	}, nil
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
