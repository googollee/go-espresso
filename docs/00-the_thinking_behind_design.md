# The Thinking Behind Design

I have a lot of personal side projects, which need backend services. When I built them, I found Go API frameworks are not good enough for personal projects. Most of frameworks belong to two kinds:

1. gRPC based.
  - API definitions and implementation in different files.
  - Require a extra step to generate scaffold files.
  - New concepts and grammar.
2. Std `net/http` based.
  - Inbound/outbond uses `any`, not type-safe.
  - Require A lot of code to parse/convert requests/responses.
  - No telemetry integration.

I hope to create a new framework for personal developers and small groups, to simplify the developing workflow but also keep best practices like type-safe, easy testing, and fast developing.

The result is `go-espresso`. This project follows guidelines below, to achieve purposes mentioned above.

- The endpoint definition and implementation are in the same place,
- No code generation,
- Provide generics helpers to reduce scaffold codes and keep type-safe.
- Integrate telemetry through middlewares.
- Follow Go guidelines, as much as possible.

Please also check other documents about design and usage.
