package espresso

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"

	"github.com/timewasted/go-accept-headers"
)

type Codec interface {
	Mime() string
	Encode(w io.Writer, v any) error
	Decode(r io.Reader, v any) error
}

type CodecJSON struct{}

func (c CodecJSON) Mime() string {
	return "applicationi/json"
}

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

func WithCodec(defaultCodec Codec, addons ...Codec) ServerOption {
	return func(s *Server) error {
		if defaultCodec != nil {
			s.codecs.defaultCodec = defaultCodec
		}

		for _, codec := range addons {
			s.codecs.all[codec.Mime()] = codec
		}

		return nil
	}
}

func (m *codecManager) decideCodec(r *http.Request) (request Codec, response Codec) {
	request = m.defaultCodec

	reqMime := r.Header.Get("Content-Type")
	if codec, ok := m.all[reqMime]; ok {
		request = codec
	}

	response = request

	acceptMime := r.Header.Get("Accept")
	accepts := accept.Parse(acceptMime)
	if len(accepts) > 0 {
		sort.Slice(accepts, func(i, j int) bool {
			return accepts[i].Q > accepts[j].Q
		})
		mime := accepts[0].Type + "/" + accepts[0].Subtype
		if codec, ok := m.all[mime]; ok {
			response = codec
		}
	}

	if codec, ok := m.all[acceptMime]; ok {
		response = codec
	}

	return
}
