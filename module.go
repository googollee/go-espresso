package espresso

import (
	"context"
	"fmt"
)

type ctxKey string

var serverKey = ctxKey(fmt.Sprintf("%T", &Server{}))

func ServerModule(ctx context.Context) *Server {
	v := ctx.Value(serverKey)
	if v == nil {
		return nil
	}

	ret, ok := v.(*Server)
	if !ok {
		return nil
	}

	return ret
}
