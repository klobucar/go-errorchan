// Package errorchan personifies Go errors as an anime character who reacts to
// your code failing. The name is wordplay that finally lands in Go: "-chan",
// the cutesy Japanese honorific, and "chan", the channel.
//
// Errors get a personality mode that changes how they are delivered, while the
// real error information stays intact and debuggable. The persona only ever
// wraps around an error; it never hides or mangles it.
//
// # Modes
//
// Three modes shape the delivery (the underlying error is identical in each):
//
//   - [ModeDere] (default): sweet, flustered, apologetic, takes the blame.
//   - [ModeTsun]: annoyed at you, blames your code, grudgingly helpful anyway.
//   - [ModeYan]: unsettlingly affectionate about your failures. Do not use in
//     production.
//
// The current global default is read with [Mode] and set with [SetMode]; both
// are safe for concurrent use.
//
// # Preserving the error
//
// Every [StyledError] keeps the original message verbatim, appended after a
// clear separator on a single line, and implements Unwrap so that [errors.Is],
// [errors.As], and [errors.Unwrap] keep working through the persona:
//
//	styled := errorchan.Wrap(err, errorchan.WithMode(errorchan.ModeTsun))
//	errors.Is(styled, io.EOF)        // still true
//	styled.Original                  // the untouched original error
//	styled.OriginalMessage           // err.Error(), captured verbatim
//
// # Surfaces
//
// The personality can be applied through several surfaces that all share one
// styling core: [Wrap] for a single error, [Styled] to restyle every error
// flowing through a channel (the pun), [Recover] to restyle a recovered panic
// inside a deferred call, and [NewSlogHandler] to style error-valued slog
// attributes.
//
// # Determinism
//
// All randomness (intro choice, kaomoji choice, stutter placement) flows
// through an injectable source. Pass [WithSeed] or [WithRand] for fully
// reproducible output, which is what the test suite relies on.
package errorchan
