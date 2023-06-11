# go-espresso

[![Go Reference](https://pkg.go.dev/badge/github.com/googollee/go-espresso.svg)](https://pkg.go.dev/github.com/googollee/go-espresso) ![CI](https://github.com/googollee/go-espresso/actions/workflows/go.yml/badge.svg)

An web/API framework.

- For individual developers and small teams.
- Code first.
  - Focus on code, instead of switching between schemas and code.
- Type safe.
  - No casting from `any`.
- Support IDE completion.
- As small dependencies as possible.
  - `httprouter`
  - `exp/slog` for logging
    - This may go to std in the future.
  - testing
    - `go-cmp`

Requirement:

- Go >= 1.20
  - Require generics.
  - `errors.Is()` supports `interface{ Unwrap() []error }`
