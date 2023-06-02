package main

import (
	"context"

	"github.com/googollee/go-espresso"
	"github.com/googollee/go-espresso/api"
	"github.com/googollee/go-espresso/sock"
)

type User struct {
	AccessKey string
	Name      string
}

type ContextData struct {
	User *User
}

type ServerService[ClientService any] struct{}

func (s *ServerService[ClientService]) Add(ctx *sock.Context[ClientService], i, j int) (int, error) {
	return i + j, nil
}

func (s *ServerService[ClientService]) OnConnect(ctx *sock.Context[ClientService]) error {
	return nil
}

func (s *ServerService[ClientService]) OnConnected(ctx *sock.Context[ClientService]) {
	echo, err := ctx.Client.Echo(ctx, "hello")
	if err != nil {
		ctx.Close()
		return
	}

	_ = echo
}

func (s *ServerService[ClientService]) OnClose(ctx *sock.Context[ClientService]) error {
	return nil
}

type ClientService interface {
	Echo(ctx context.Context, str string) (string, error)
}

func main() {
	server := espresso.NewServer(ContextData{})

	service := ServerService{}

	sock.Handle[ClientService]("/sock", server, service).
		WithDefaultCodec(api.CodecJSON)

	server.ListenAndServe(":8080")
}
