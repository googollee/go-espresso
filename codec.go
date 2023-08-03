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

func WithCodec(defaultCodec Codec, all map[string]Codec) ServerOption {
	return func(s *Server) error {
		if defaultCodec != nil {
			s.codecs.defaultCodec = defaultCodec
		}

		s.codecs.appendCodecs(all)

		return nil
	}
}

func (m *codecManager) appendCodecs(all map[string]Codec) {
	if m.all == nil {
		m.all = make(map[string]Codec)
	}

	for mime, codec := range all {
		m.all[mime] = codec
	}
}

func (m *codecManager) decideCodec(r *http.Request) Codec {
	if m.defaultCodec != nil {
		return m.defaultCodec
	}

	return CodecJSON{}
}
