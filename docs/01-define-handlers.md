# Define Handlers

The signature of a `espresso` handler is:

```go
type Handler func(ctx espresso.Context) error
```

It could be implemented by either functions or struct methods:

```go
func FuncHanlder(ctx espresso.Context) error {
    // ...
}

type Service struct{
    // ...
}

func (s *Service) MethodHandlerWithPointer(ctx espresso.Context) error {
    // ...
}

func (s Service) MethodHandlerWithValue(ctx espresso.Context) error {
    // ...
}
```

To be simple, examples below are based on function handlers.

Use the `Context` instance in a handler to define endpoint information, like below:

```go
func Handler(ctx espresso.Context) error {
    var param int
    if err := ctx.Endpoint(http.MethodGet, "/endpoint/with/:param_in_path").
        BindPath("param_in_path", &param).
        End(); err != nil {
        return err
    }
    // ...
}
```

`Context.Endpoint()` registers an endpoint as `GET /endpoint/with/:param_in_path`, with binding values:

- binds the parameter `:param_in_path` in a path to an integer `param` variable.
  - Valid paths are like `/endpoint/with/1` or `/endpoint/with/1000`, but not `/endpoint/with/non_number`.
  - Parse `:param_in_path` part in the path of a request to a `int` value and assign to the `param` variable.

When registering this handler, `espresso` passes a special `Context` to collect bindings with `Context.Endpoint()`, and panic in `End()`. Please put this code block at the top of a handler, to avoid calling real logic code below when registering.

When handling a request, `espresso` passes another `Context` to this handler, parse values from strings in the request and assign results to bind variables. If there are parsing errors, all errors return by `End()`.
