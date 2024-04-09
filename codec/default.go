package codec

import "encoding/json"

type JSON struct{}

func (JSON) Mime() string {
	return "application/json"
}

func (JSON) DecodeRequest(ctx Context, v any) error {
	return json.NewDecoder(ctx.Request().Body).Decode(v)
}

func (JSON) EncodeResponse(ctx Context, v any) error {
	return json.NewEncoder(ctx.ResponseWriter()).Encode(v)
}
