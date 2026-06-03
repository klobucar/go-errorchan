# go-errorchan

> A pun that finally compiles in Go: **-chan**, the cutesy anime honorific, *and* **chan**, the channel your errors flow through. `go-errorchan` personifies your errors as an anime character who reacts to your code failing — while keeping the real error completely intact and debuggable.

[![CI](https://github.com/klobucar/go-errorchan/actions/workflows/ci.yml/badge.svg)](https://github.com/klobucar/go-errorchan/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/klobucar/go-errorchan.svg)](https://pkg.go.dev/github.com/klobucar/go-errorchan)
[![Go Report Card](https://goreportcard.com/badge/github.com/klobucar/go-errorchan)](https://goreportcard.com/report/github.com/klobucar/go-errorchan)

```go
import "github.com/klobucar/go-errorchan"
```

The headline feature is `Styled`, which restyles **every error flowing through a channel** — because in Go, `chan` was a channel all along.

## Install

```sh
go get github.com/klobucar/go-errorchan
```

Standard library only. Zero third-party runtime dependencies.

## Modes

Three personality modes change how an error is *delivered*. The underlying error is byte-for-byte identical in every mode.

| Mode | Vibe | Sample |
| --- | --- | --- |
| `dere` (default) | sweet, flustered, apologetic, takes the blame | `evewything was finye u-untiw this *errors.errorString*... g-gomen nyasai >_< — connection refused` |
| `tsun` | annoyed at you, blames your code, grudgingly helpful | `tch, a *errors.errorString*. finye, hewe's the detaiw. don't make me wepeat it >:( — EOF` |
| `yan` | unsettlingly affectionate about your failures (**do not use in prod**) | `anyothew *errors.errorString*... good. t-the mowe you faiw, the mowe you'we minye (◡‿◡✿) — permission denied` |

Each mode has its own intro lines and kaomoji, run through a shared phonetic transform (`r`/`l` → `w`, `n`+vowel → `ny`, and seed-driven stutters for the heavier modes).

```go
errorchan.SetMode("tsun")  // global default; concurrency-safe
mode := errorchan.Mode()   // "tsun"
```

## The one rule

**The real error is never hidden or mangled.** The persona only wraps around it:

- `errors.Is` and `errors.As` see straight through (`StyledError` implements `Unwrap`).
- The original is exposed on the value: `styled.Original` (the error) and `styled.OriginalMessage` (its message, verbatim).
- `Error()` keeps the original message verbatim after a clear ` — ` separator, always on a single line so it never wrecks your logs.
- Re-wrapping preserves `%w` semantics, so callers' sentinel and type checks keep passing.

```go
styled := errorchan.Wrap(io.EOF, errorchan.WithMode("tsun"))
errors.Is(styled, io.EOF)   // true
styled.OriginalMessage      // "EOF"
```

## API

```go
// Raw phonetic transform, exposed.
errorchan.Uwuify("really cool")                       // "weawwy coow"

// Wrap a single error: *StyledError, Unwrap() == err.
styled := errorchan.Wrap(err, opts...)

// THE PUN: restyle every error flowing through a channel.
out := errorchan.Styled(in <-chan error, opts...)     // returns <-chan error

// Restyle a recovered panic into a named error return.
func doThing() (err error) {
    defer errorchan.Recover(&err, opts...)
    // ...
}

// Style error-valued attributes in slog records.
h := errorchan.NewSlogHandler(base, opts...)
log := slog.New(h)
```

### Options

Configuration uses the functional-options pattern:

- `WithMode(string)` — override the mode for this operation.
- `WithSeed(int64)` — deterministic output from a fresh, per-operation source (safe to reuse across goroutines).
- `WithRand(*rand.Rand)` — deterministic output from a source you own (not safe to share across goroutines).

All randomness flows through the injected source, so seeded output is fully reproducible — which is exactly how the test suite pins every example.

```go
errorchan.Uwuify("really cool", errorchan.WithSeed(1)) // always "weawwy coow"
```

## The channel surface

```go
in := make(chan error, 2)
in <- errors.New("timeout")
in <- errors.New("refused")
close(in)

for err := range errorchan.Styled(in, errorchan.WithMode("dere")) {
    fmt.Println(err)
}
```

`Styled` resolves its config once and shares a single source across the stream, so output varies error-to-error yet stays reproducible under a seed. `nil` values pass through as `nil` (never a non-nil wrapper around a nil error).

## Design

The package is built in layers, each independent of the ones above it:

1. **Phonetic transform** (`transform.go`) — deterministic given an injected `*rand.Rand`, with no knowledge of errors or personas.
2. **Persona layer** (`persona.go`) — unexported data mapping each mode to its framing lines and kaomoji.
3. **Usage surfaces** (`Wrap`, `Styled`, `Recover`, `NewSlogHandler`) — thin shells over one shared styling core.

## Development

```sh
go test -race -cover ./...   # ~100% coverage, race-clean
gofmt -l .                   # formatting
go vet ./...                 # vet
golangci-lint run ./...      # strict lint (see .golangci.yml)
govulncheck ./...            # vulnerability scan
```

CI runs all of the above across a Go version matrix on every push and PR.

## License

[MIT](LICENSE)
