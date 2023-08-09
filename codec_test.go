package espresso

import (
	"io"
	"net/http"
	"testing"
)

type fakeCodec struct {
	mime string
}

func (c *fakeCodec) Mime() string {
	return c.mime
}

func (c *fakeCodec) Encode(w io.Writer, v any) error {
	return nil
}

func (c *fakeCodec) Decode(r io.Reader, v any) error {
	return nil
}

func TestCodecDefaultManager(t *testing.T) {
	tests := []struct {
		name         string
		headers      map[string]string
		wantReqMime  string
		wantRespMime string
	}{
		{
			name: "RequestWithXML",
			headers: map[string]string{
				"Content-Type": "application/xml",
			},
			wantReqMime:  "application/json",
			wantRespMime: "application/json",
		},
		{
			name: "AllXML",
			headers: map[string]string{
				"Content-Type": "application/xml",
				"Accept":       "application/xml",
			},
			wantReqMime:  "application/json",
			wantRespMime: "application/json",
		},
		{
			name:         "NoHeader",
			wantReqMime:  "application/json",
			wantRespMime: "application/json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.NewRequest(http.MethodPost, "/", nil)
			if err != nil {
				t.Fatalf("Create request error: %v", err)
			}

			for key, value := range tc.headers {
				r.Header.Add(key, value)
			}

			cm := defaultManager()
			reqC, respC := cm.decideCodec(r)
			if got, want := reqC.Mime(), tc.wantReqMime; got != want {
				t.Errorf("defaultManager().decideCodec()[0].Mime() = %q, want: %q", got, want)
			}
			if got, want := respC.Mime(), tc.wantRespMime; got != want {
				t.Errorf("defaultManager().decideCodec()[1].Mime() = %q, want: %q", got, want)
			}
		})
	}
}
func TestCodecMore(t *testing.T) {
	xmlCodec := &fakeCodec{mime: "application/xml"}
	yamlCodec := &fakeCodec{mime: "application/yaml"}
	tests := []struct {
		name         string
		headers      map[string]string
		wantReqMime  string
		wantRespMime string
	}{
		{
			name:         "NoHeader",
			wantReqMime:  "application/json",
			wantRespMime: "application/json",
		},
		{
			name: "RequestWithXML",
			headers: map[string]string{
				"Content-Type": "application/xml",
			},
			wantReqMime:  "application/xml",
			wantRespMime: "application/xml",
		},
		{
			name: "ResponseWithXML",
			headers: map[string]string{
				"Accept": "application/xml",
			},
			wantReqMime:  "application/json",
			wantRespMime: "application/xml",
		},
		{
			name: "AllYAML",
			headers: map[string]string{
				"Content-Type": "application/yaml",
				"Accept":       "application/yaml",
			},
			wantReqMime:  "application/yaml",
			wantRespMime: "application/yaml",
		},
		{
			name: "MultipleResponse1st",
			headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/xml;q=1, application/json;q=0.8",
			},
			wantReqMime:  "application/json",
			wantRespMime: "application/xml",
		},
		{
			name: "MultipleResponse2nd",
			headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/nonexist;q=1, application/xml;q=0.8",
			},
			wantReqMime:  "application/json",
			wantRespMime: "application/xml",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.NewRequest(http.MethodPost, "/", nil)
			if err != nil {
				t.Fatalf("Create request error: %v", err)
			}

			for key, value := range tc.headers {
				r.Header.Add(key, value)
			}

			svr, err := New(WithCodec(CodecJSON{}, xmlCodec, yamlCodec))
			if err != nil {
				t.Fatalf("WithCodec() error: %v", err)
			}
			cm := svr.codecs

			reqC, respC := cm.decideCodec(r)
			if got, want := reqC.Mime(), tc.wantReqMime; got != want {
				t.Errorf("defaultManager().decideCodec()[0].Mime() = %q, want: %q", got, want)
			}
			if got, want := respC.Mime(), tc.wantRespMime; got != want {
				t.Errorf("defaultManager().decideCodec()[1].Mime() = %q, want: %q", got, want)
			}
		})
	}
}
