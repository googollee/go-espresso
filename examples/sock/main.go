//go:build wip

package main

import (
	"context"

	"github.com/googollee/go-espresso"
)

type User struct {
	AccessKey string
	Name      string
}

type ContextData struct {
	User *User
}

type ServerService struct{}

func (s *ServerService) Add(ctx sock.Context[ContextData, ClientService], i, j int) (int, error) {
	return i + j, nil
}

func (s *ServerService) OnConnect(ctx sock.Context[ContextData, ClientService]) error {
	return nil
}

func (s *ServerService) OnConnected(ctx sock.Context[ContextData, ClientService]) {
	echo, err := ctx.Client.Echo(ctx, "hello")
	if err != nil {
		ctx.Close()
		return
	}

	_ = echo
}

func (s *ServerService) OnClose(ctx sock.Context[ClientService, ClientService]) error {
	return nil
}

type ClientService interface {
	Echo(ctx context.Context, str string) (string, error)
}

func main() {
	server := espresso.NewServer()

	service := ServerService{}

	sock.Handle[ClientService]("/sock", server, service).
		WithDefaultCodec(api.CodecJSON)

	server.ListenAndServe(":8080")
}
