# Contributing to go-errorchan

Thanks for your interest in contributing! This is a small, dependency-free Go
library, and contributions of all sizes are welcome — bug fixes, new personas,
docs, and tests.

## Getting started

```sh
git clone https://github.com/klobucar/go-errorchan
cd go-errorchan
go test ./...
```

The module targets the Go version pinned in [`go.mod`](go.mod) and uses only the
standard library — please keep it that way (zero third-party runtime
dependencies).

## Before you open a pull request

Run the same checks CI runs:

```sh
gofmt -l .                   # must print nothing
go vet ./...
go test -race -cover ./...   # coverage is gated at 95%+ (currently ~100%)
golangci-lint run ./...      # strict lint; see .golangci.yml
govulncheck ./...            # vulnerability scan
```

A change is ready when all of the above pass locally. CI runs them across a Go
version matrix on every push and PR.

## Guidelines

- **The one rule:** the real error is never hidden or mangled. Any new surface
  must keep the original error reachable via `Unwrap` (so `errors.Is` /
  `errors.As` see through it) and preserve `OriginalMessage` verbatim. See the
  README's "The one rule" section.
- **Determinism:** all randomness must flow through the injected `*rand.Rand` so
  output stays reproducible under `WithSeed` / `WithRand`. Tests pin behavior
  with seeds — please add seeded tests for new behavior.
- **Personas:** new intro lines and kaomoji go in `persona.go`. Intros are plain
  English (the phonetic transform mangles them at delivery time) and may use a
  single `%s` for the error's Go type.
- Keep the public API small and documented; exported symbols need doc comments
  that render well on [pkg.go.dev](https://pkg.go.dev/github.com/klobucar/go-errorchan).

## Reporting bugs and requesting features

Open an issue using one of the templates. For bugs, a small reproducer (ideally
with a `WithSeed` so it's deterministic) is the fastest path to a fix.

## License

By contributing, you agree that your contributions are licensed under the
project's [MIT License](LICENSE).
