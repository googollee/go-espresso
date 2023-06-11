package espresso_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/googollee/go-espresso"
)

func getUserFromSessionID(id string) string {
	return ""
}

func ExampleHandler_middlewareAuth() {
	type Data struct {
		User string
	}

	getUser := func(ctx espresso.Context[Data]) error {
		cookie, err := ctx.Request().Cookie("session_id")
		if err != nil {
			return espresso.ErrWithStatus(http.StatusUnauthorized, err)
		}
		sessionID := cookie.Value

		// this function is provided by your system.
		ctx.Data().User = getUserFromSessionID(sessionID)

		return nil
	}

	_ = getUser
}

func ExampleHandler_middlewareWithTimeout() {
	type Data struct {
		// Your data
	}

	withTimeout := func(timeout time.Duration) espresso.Handler[Data] {
		return func(ctx espresso.Context[Data]) error {
			cancelCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			newCtx := ctx.WithContext(cancelCtx)
			newCtx.Next()

			return nil
		}
	}

	_ = withTimeout(time.Minute)
}

func ExampleHandler_middlewareErrorResponse() {
	type Data struct {
		// Your data
	}

	errorResponse := func(ctx espresso.Context[Data]) error {
		ctx.Next()

		if ctx.Err() == nil {
			return nil
		}

		code := http.StatusInternalServerError
		if he, ok := ctx.Err().(espresso.HTTPCoder); ok {
			code = he.HTTPCode()
		}

		ctx.ResponseWriter().Header().Add("Context-Type", "text/plain")
		ctx.ResponseWriter().WriteHeader(code)
		fmt.Fprintf(ctx.ResponseWriter(), "with error: %v", ctx.Err())

		return nil
	}

	_ = errorResponse
}
