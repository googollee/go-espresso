# Define Endpoints

Use the `Context` instance in a handler to define an endpoint, like below:

```go
func EndpointHandler(ctx espresso.Context) error {
    var param int
    if err := ctx.Endpoint(http.MethodGet, "/endpoint/with/:param_in_path").
        BindPath("param_in_path", &param).
        End(); err != nil {
        return err
    }
    // ...
}
```

`Context.Endpoint()` registers an endpoint as `GET /endpoint/with/:param_in_path`, and binds the parameter `:param_in_path` in a path to the `param` variable. The `param` is `int` type, so valid paths are like `/endpoint/with/1` or `/endpoint/with/1000`, but not `/endpoint/with/non_number`. When registering this handler, `go-espresso` passes a special `Context` to collect data with `Context.Endpoint()`, and panic in `End()`. Please put this code block at the top of a handler, to avoid calling real logic code below when registering. When handling a request, `go-espresso` passes another `Context` to this handler, and parse `:param_in_path` part in the path to a `int` value and assign to the `param` variable. If there are parsing errors, all errors return by `End()`.

You can define a struct with many methods as handlers, and register them at once:

```go
service := &Service{}

server := espresso.New()
server.HandleAll(service)
```
