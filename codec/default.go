package codec

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

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

type YAML struct{}

func (YAML) Mime() string {
	return "application/yaml"
}

func (YAML) DecodeRequest(ctx Context, v any) error {
	return yaml.NewDecoder(ctx.Request().Body).Decode(v)
}

func (YAML) EncodeResponse(ctx Context, v any) error {
	encoder := yaml.NewEncoder(ctx.ResponseWriter())
	defer encoder.Close()

	return encoder.Encode(v)
}
