package codec

import (
	"context"
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"
)

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
