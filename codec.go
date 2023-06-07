package espresso

import (
	"encoding/json"
	"io"
)

type Encoder interface {
	Encode(v any) error
}

type Decoder interface {
	Decode(v any) error
}

type Codec interface {
	Mime() string
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

var (
	CodecJSON Codec = codecJSON{}
)

type codecJSON struct{}

func (c codecJSON) Mime() string {
	return "application/json"
}

func (c codecJSON) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (c codecJSON) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}
