package espresso

import (
	"encoding/json"
	"io"
)

type Codec interface {
	Encode(w io.Writer, v any) error
	Decode(r io.Reader, v any) error
}

var DefaultCodec Codec = CodecJSON{}

type CodecJSON struct{}

func (c CodecJSON) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func (c CodecJSON) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
