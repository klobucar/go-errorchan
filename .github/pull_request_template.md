<!-- Thanks for contributing to go-errorchan! -->

## What & why

<!-- What does this change do, and why? Link any related issue: Fixes #123 -->

## Checklist

- [ ] `gofmt -l .` prints nothing
- [ ] `go vet ./...` passes
- [ ] `go test -race -cover ./...` passes (coverage stays at 95%+)
- [ ] `golangci-lint run ./...` passes
- [ ] The wrapped error is still reachable (`errors.Is`/`errors.As` see through; `OriginalMessage` verbatim)
- [ ] New behavior has seeded, deterministic tests
- [ ] Exported symbols have doc comments
