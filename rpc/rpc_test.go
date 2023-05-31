package rpc

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/googollee/go-espresso"
)

type Arg struct {
	I int `json:"i"`
}

type Reply struct {
	Double int `json:"double"`
}

type Data struct{}

func TestRPC(t *testing.T) {
	svr := espresso.NewServer(Data{})
	svc, err := New(svr, "/prefix")
	if err != nil {
		t.Fatalf("create service error: %v", err)
	}

	POST(svc, "/add", func(ctx *espresso.Context[Data], arg *Arg) (*Reply, error) {
		return &Reply{
			Double: arg.I * 2,
		}, nil
	})

	testSvc := httptest.NewServer(svr)
	defer testSvc.Close()

	client := testSvc.Client()

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(Arg{
		I: 10,
	})

	resp, err := client.Post(testSvc.URL+"/prefix/add", "application/json", &buf)
	if err != nil {
		t.Fatalf("post to /add error: %v", err)
	}
	defer resp.Body.Close()

	var reply Reply
	if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		t.Fatalf("decode response from /add error: %v", err)
	}

	if got, want := reply.Double, 20; got != want {
		t.Errorf("response is %d, want %d", got, want)
	}
}
