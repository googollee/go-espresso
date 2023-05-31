package rpc

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

type CodecJSON struct{}

func (c CodecJSON) Mime() string {
	return "application/json"
}

func (c CodecJSON) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (c CodecJSON) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}
