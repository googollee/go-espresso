package espresso

import (
	"encoding/json"
	"io"
	"net/http"
)

type Codec interface {
	Encode(w io.Writer, v any) error
	Decode(r io.Reader, v any) error
}

type CodecJSON struct{}

func (c CodecJSON) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func (c CodecJSON) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

type codecManager struct {
	defaultCodec Codec
	all          map[string]Codec
}

func defaultManager() codecManager {
	return codecManager{
		defaultCodec: CodecJSON{},
		all:          make(map[string]Codec),
	}
}

func WithCodec(defaultCodec Codec, all map[string]Codec) ServerOption {
	return func(s *Server) error {
		if defaultCodec != nil {
			s.codecs.defaultCodec = defaultCodec
		}

		for mime, codec := range all {
			s.codecs.all[mime] = codec
		}

		return nil
	}
}

func (m *codecManager) decideCodec(r *http.Request) Codec {
	if m.defaultCodec != nil {
		return m.defaultCodec
	}

	return CodecJSON{}
}
