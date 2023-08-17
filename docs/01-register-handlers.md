# Register Handlers

It's easy to register a function handler:

```go
func FuncHandler(ctx espresso.Context) error {
    // ...
}

svr := espresso.New()

svr.HandleFunc(FuncHandler)
```

It's also easy to register all method handlers in a struct:

```go
type Service struct{
    // ...
}

func (s *Service) Handle1(ctx espresso.Context) error { /* ... */ }
func (s *Service) Handle2(ctx espresso.Context) error { /* ... */ }
func (s *Service) Handle3(ctx espresso.Context) error { /* ... */ }
func (s *Service) Handle4(ctx espresso.Context) error { /* ... */ }

service := &Service{}

svr := espresso.New()
svr.HandleAll(service)
```

`HandleAll` goes through all methods of the given value by reflecting, and register any methods matching with the signature `func(espresso.Context) error` of the `espresso` handler. When handling requests, it calls methods directly. No reflecting during handling real requests.
